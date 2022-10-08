package http

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lbryio/lbry.go/v2/extras/errors"
)

func getOptimizationParams(c *gin.Context) (width int64, height int64, quality int64, err error) {
	width, height, err = getDimensions(c)
	if err != nil {
		return width, height, 0, err
	}
	quality, err = strconv.ParseInt(strings.TrimPrefix(c.Param("quality"), ":"), 10, 32)
	if err != nil {
		return width, height, quality, errors.Err(err)
	}
	return width, height, quality, nil
}

func getDimensions(c *gin.Context) (width int64, height int64, err error) {
	dimensions := strings.Split(c.Param("dimensions"), ":")
	if len(dimensions) != 3 {
		return 0, 0, errors.Err("dimensions should be in the form of /s:width:height/")
	}
	width, err = strconv.ParseInt(dimensions[1], 10, 32)
	if err != nil {
		return 0, 0, errors.Err(err)
	}
	height, err = strconv.ParseInt(dimensions[2], 10, 32)
	if err != nil {
		return width, 0, errors.Err(err)
	}
	//the frontend is requesting something it doesn't actually want.... this forces us to hardcode it here
	//TODO: get rid of this and have the frontend NOT pass both params
	//this will also mess up the caches as things will be cached with improper parameters
	if width != 0 && height != 0 {
		height = 0
	}
	return width, height, nil
}

func handleExceptions(c *gin.Context, width int64, height int64, quality int64, path string) (redirected bool) {
	urlToProxy := extractUrl(c)
	imgurUrl := regexp.MustCompile(`^https?://i?\.?imgur\.com/.+?$`)
	// temporarily disable imgur proxying because of throttling
	if imgurUrl.MatchString(urlToProxy) {
		c.Redirect(http.StatusTemporaryRedirect, urlToProxy)
		return true
	}
	malformedSpeechUrl := strings.Index(urlToProxy, "https://spee.ch/") == 0
	if malformedSpeechUrl {
		urlToProxy = strings.TrimPrefix(urlToProxy, "https://spee.ch/")
		if parts := regexp.MustCompile(`^(view/)?([a-f0-9]+)/(.*?)\.(.*)$`).FindStringSubmatch(urlToProxy); parts != nil {
			c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf(path, width, height, quality, url.QueryEscape(fmt.Sprintf("https://player.odycdn.com/speech/%s:%s.%s", parts[3], parts[2], parts[4]))))
			return true
		}
	}
	atWebp := strings.HasSuffix(urlToProxy, "@webp")
	if atWebp {
		urlToProxy = strings.TrimSuffix(urlToProxy, "@webp")
		c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf(path, width, height, quality, url.QueryEscape(urlToProxy)))
		return true
	}
	oldSpeechBug := strings.HasSuffix(urlToProxy, "..jpeg") || strings.HasSuffix(urlToProxy, "..png")
	if oldSpeechBug {
		urlToProxy = strings.TrimSuffix(urlToProxy, "..jpeg")
		urlToProxy = strings.TrimSuffix(urlToProxy, "..png")
		c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf(path, width, height, quality, url.QueryEscape(urlToProxy)))
		return true
	}
	decommissionedProxy := strings.Contains(urlToProxy, "https://lbry-boost.org/redirect-event?source=")
	if decommissionedProxy {
		urlToProxy = strings.Replace(urlToProxy, "https://lbry-boost.org/redirect-event?source=", "", -1)
		c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf(path, width, height, quality, url.QueryEscape(urlToProxy)))
		return true
	}
	recursionUrlToProxy := strings.TrimPrefix(c.Param("url"), "/")
	hasRecursion := strings.Contains(recursionUrlToProxy, "https://thumbnails.odycdn.com")
	if hasRecursion {
		cutIndex := strings.LastIndex(recursionUrlToProxy, "plain/") + 6
		if cutIndex > 6 {
			urlToProxy = recursionUrlToProxy[cutIndex:len(recursionUrlToProxy)]
			c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf(path, width, height, quality, url.QueryEscape(urlToProxy)))
			return true
		}
		_ = c.AbortWithError(http.StatusBadRequest, errors.Err("malformed recursive URL"))
		return true
	}
	return false
}

func extractUrl(c *gin.Context) string {
	urlToProxy := strings.TrimPrefix(c.Param("url"), "/")
	uriSplit := strings.Split(c.Request.RequestURI, urlToProxy)
	queryString := ""
	if len(uriSplit) > 1 {
		queryString = uriSplit[1]
	}
	urlToProxy += queryString
	return urlToProxy
}
