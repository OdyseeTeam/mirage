package downloader

import (
	"io"
	"net/http"
	"time"

	"github.com/lbryio/lbry.go/v2/extras/errors"
)

func DownloadFile(URL string, isRetry bool) ([]byte, error) {
	method := "GET"

	client := &http.Client{
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest(method, URL, nil)
	if err != nil {
		return nil, errors.Err(err)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.Err(err)
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		if !isRetry && response.StatusCode == http.StatusBadGateway {
			time.Sleep(100 * time.Millisecond)
			return DownloadFile(URL, true)
		}
		return nil, errors.Err("Received non 200 response code %d for %s", response.StatusCode, URL)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Err(err)
	}

	return bodyBytes, nil
}
