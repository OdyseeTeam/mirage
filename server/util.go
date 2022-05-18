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

	return width, height, nil
}

func handleExceptions(c *gin.Context, width int64, height int64, quality int64) {
	urlToProxy := extractUrl(c)
	malformedSpeechUrl := strings.Index(urlToProxy, "https://spee.ch/") == 0
	if malformedSpeechUrl {
		urlToProxy = strings.TrimPrefix(urlToProxy, "https://spee.ch/")
		if parts := regexp.MustCompile(`^(view/)?([a-f0-9]+)/(.*?)\.(.*)$`).FindStringSubmatch(urlToProxy); parts != nil {
			c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/optimize/s:%d:%d/quality:%d/plain/%s", width, height, quality, url.QueryEscape(fmt.Sprintf("https://player.odycdn.com/speech/%s:%s.%s", parts[3], parts[2], parts[4]))))
			return
		}
	}
	atWebp := strings.HasSuffix(urlToProxy, "@webp")
	if atWebp {
		urlToProxy = strings.TrimSuffix(urlToProxy, "@webp")
		c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("/optimize/s:%d:%d/quality:%d/plain/%s", width, height, quality, url.QueryEscape(urlToProxy)))
		return
	}
	oldSpeechBug := strings.HasSuffix(urlToProxy, "..jpeg")
	if oldSpeechBug {
		urlToProxy = strings.TrimSuffix(urlToProxy, "..jpeg")
		c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("/optimize/s:%d:%d/quality:%d/plain/%s", width, height, quality, url.QueryEscape(urlToProxy)))
		return
	}
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
