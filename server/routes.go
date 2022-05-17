package http

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/OdyseeTeam/mirage/downloader"

	"github.com/OdyseeTeam/gody-cdn/store"
	"github.com/gin-gonic/gin"
	"github.com/golang/groupcache/singleflight"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/sirupsen/logrus"
)

type optimizerParams struct {
	Width      int64  `json:"width"`
	Height     int64  `json:"height"`
	Quality    int64  `json:"quality"`
	UrlToProxy string `json:"urlToProxy"`
}

var sf = singleflight.Group{}

func (s *Server) optimizeHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered from panic: %v", r)
		}
	}()
	dimensions := strings.Split(c.Param("dimensions"), ":")
	if len(dimensions) != 3 {
		_ = c.AbortWithError(http.StatusBadRequest, errors.Err("dimensions should be in the form of /s:width:height/"))
	}
	width, err := strconv.ParseInt(dimensions[1], 10, 32)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.Err(err))
		return
	}
	height, err := strconv.ParseInt(dimensions[2], 10, 32)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.Err(err))
		return
	}
	quality, err := strconv.ParseInt(strings.TrimPrefix(c.Param("quality"), ":"), 10, 32)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.Err(err))
		return
	}
	urlToProxy := strings.TrimPrefix(c.Param("url"), "/")
	var oP = optimizerParams{
		Width:      width,
		Height:     height,
		Quality:    quality,
		UrlToProxy: urlToProxy,
	}
	type OptimizedImage struct {
		optimizedImage   []byte
		originalMimeType string
		originalSize     int
		optimizedSize    int
	}

	key := fmt.Sprintf("%s-%d-%d-%d", urlToProxy, width, height, quality)
	v, err := sf.Do(key, func() (interface{}, error) {
		h := sha1.New()
		h.Write([]byte(key))
		hashedName := hex.EncodeToString(h.Sum(nil))
		obj, _, err := s.cache.Get(hashedName, nil)
		if err == nil {
			return OptimizedImage{
				optimizedImage:   obj,
				originalMimeType: "unknown",
				originalSize:     0,
				optimizedSize:    len(obj),
			}, nil
		}
		if err != nil && !strings.Contains(err.Error(), store.ErrObjectNotFound.Error()) {
			return nil, err
		}
		image, err := downloader.DownloadFile(urlToProxy)
		if err != nil {
			return nil, err
		}
		optimized, origMime, err := s.optimizer.Optimize(image, quality)
		if err != nil {
			logrus.Errorf("failed to optimize resource with content type: %s", origMime)
			return nil, err
		}
		err = s.cache.Put(hashedName, optimized, nil)
		if err != nil {
			logrus.Errorf("error storing %s: %s", key, errors.FullTrace(err))
		}
		return OptimizedImage{
			optimizedImage:   optimized,
			originalMimeType: origMime,
			originalSize:     len(image),
			optimizedSize:    len(optimized),
		}, nil
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.Err(err))
		return
	}
	optimizedData, ok := v.(OptimizedImage)
	if !ok {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.Err("could not cast from sf cache"))
		return
	}

	c.Header("Content-Length", fmt.Sprintf("%d", optimizedData.optimizedSize))
	c.Header("X-mirage-saved-bytes", fmt.Sprintf("%d", optimizedData.originalSize-optimizedData.optimizedSize))
	c.Header("X-mirage-compression-ratio", fmt.Sprintf("%.2f:1", float64(optimizedData.originalSize)/float64(optimizedData.optimizedSize)))
	c.Header("X-mirage-original-mime", optimizedData.originalMimeType)
	c.Data(200, "image/webp", optimizedData.optimizedImage)
	logrus.Infof("%v", oP)
}

func (s *Server) recoveryHandler(c *gin.Context, err interface{}) {
	c.JSON(500, gin.H{
		"title": "Error",
		"err":   err,
	})
}

func (s *Server) ErrorHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}
		logrus.Errorln(errors.FullTrace(err))
		c.String(-1, err.Error())
		return
	}
}

func (s *Server) addCSPHeaders(c *gin.Context) {
	c.Header("Report-To", `{"group":"default","max_age":31536000,"endpoints":[{"url":"https://6fd448c230d0731192f779791c8e45c3.report-uri.com/a/d/g"}],"include_subdomains":true}`)
	c.Header("Content-Security-Policy", "script-src 'none'; report-uri https://6fd448c230d0731192f779791c8e45c3.report-uri.com/r/d/csp/enforce; report-to default")
}
