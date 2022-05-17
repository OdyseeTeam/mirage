package http

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/OdyseeTeam/mirage/downloader"
	"github.com/OdyseeTeam/mirage/metadata"

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

func (s *Server) simpleRedirect(c *gin.Context) {
	urlToProxy := strings.TrimPrefix(c.Param("url"), "/")
	uriSplit := strings.Split(c.Request.RequestURI, urlToProxy)
	queryString := ""
	if len(uriSplit) > 1 {
		queryString = uriSplit[1]
	}
	urlToProxy += queryString
	c.Redirect(http.StatusTemporaryRedirect, "/optimize/s:0:0/quality:85/plain/"+url.QueryEscape(urlToProxy))
}
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
	uriSplit := strings.Split(c.Request.RequestURI, urlToProxy)
	queryString := ""
	if len(uriSplit) > 1 {
		queryString = uriSplit[1]
	}
	urlToProxy += queryString

	type OptimizedImage struct {
		optimizedImage *[]byte
		metadata       *metadata.ImageMetadata
	}

	key := fmt.Sprintf("%s-%d-%d-%d", urlToProxy, width, height, quality)
	v, err := sf.Do(key, func() (interface{}, error) {
		h := sha1.New()
		h.Write([]byte(key))
		hashedName := hex.EncodeToString(h.Sum(nil))
		obj, _, err := s.cache.Get(hashedName, nil)
		if err == nil {
			md, err := s.metadataManager.Retrieve(hashedName)
			if md == nil {
				if err != nil {
					logrus.Errorf("cannot retrieve metadata: %s", errors.FullTrace(err))
				}
				md = &metadata.ImageMetadata{
					OriginalURL:      urlToProxy,
					GodycdnHash:      hashedName,
					Checksum:         fmt.Sprintf("%x", sha256.Sum256(obj)),
					OriginalMimeType: "unknown",
					OriginalSize:     0,
					OptimizedSize:    len(obj),
				}
			}
			return OptimizedImage{
				optimizedImage: &obj,
				metadata:       md,
			}, nil
		}
		if err != nil && !strings.Contains(err.Error(), store.ErrObjectNotFound.Error()) {
			return nil, err
		}
		image, err := downloader.DownloadFile(urlToProxy)
		if err != nil {
			return nil, err
		}
		optimized, origMime, err := s.optimizer.Optimize(image, quality, width, height)
		if err != nil {
			logrus.Errorf("failed to optimize resource with content type: %s", origMime)
			return nil, err
		}
		err = s.cache.Put(hashedName, optimized, nil)
		if err != nil {
			logrus.Errorf("error storing %s: %s", key, errors.FullTrace(err))
		}
		md := &metadata.ImageMetadata{
			OriginalURL:      urlToProxy,
			GodycdnHash:      hashedName,
			Checksum:         fmt.Sprintf("%x", sha256.Sum256(optimized)),
			OriginalMimeType: origMime,
			OriginalSize:     len(image),
			OptimizedSize:    len(optimized),
		}
		err = s.metadataManager.Persist(md)
		if err != nil {
			logrus.Errorf("failed to persiste metadata for object %s: %s", urlToProxy, errors.FullTrace(err))
		}
		return OptimizedImage{
			optimizedImage: &optimized,
			metadata:       md,
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

	c.Header("Content-Length", fmt.Sprintf("%d", optimizedData.metadata.OptimizedSize))
	c.Header("X-mirage-saved-bytes", fmt.Sprintf("%d", optimizedData.metadata.OriginalSize-optimizedData.metadata.OptimizedSize))
	c.Header("X-mirage-compression-ratio", fmt.Sprintf("%.2f:1", float64(optimizedData.metadata.OriginalSize)/float64(optimizedData.metadata.OptimizedSize)))
	c.Header("X-mirage-original-mime", optimizedData.metadata.OriginalMimeType)
	c.Header("Cache-control", "max-age=604800")
	c.Data(200, "image/webp", *optimizedData.optimizedImage)
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
		c.Header("Cache-control", "max-age=240")
		return
	}
}

func (s *Server) addCSPHeaders(c *gin.Context) {
	c.Header("Report-To", `{"group":"default","max_age":31536000,"endpoints":[{"url":"https://6fd448c230d0731192f779791c8e45c3.report-uri.com/a/d/g"}],"include_subdomains":true}`)
	c.Header("Content-Security-Policy", "script-src 'none'; report-uri https://6fd448c230d0731192f779791c8e45c3.report-uri.com/r/d/csp/enforce; report-to default")
}
