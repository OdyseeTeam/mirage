package downloader

import (
	"io"
	"net/http"

	"github.com/lbryio/lbry.go/v2/extras/errors"
)

func DownloadFile(URL string) ([]byte, error) {
	response, err := http.Get(URL)
	if err != nil {
		return nil, errors.Err(err)
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(response.Body)

	if response.StatusCode != 200 {
		return nil, errors.Err("Received non 200 response code %d for %s", response.StatusCode, URL)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Err(err)
	}

	return bodyBytes, nil
}
