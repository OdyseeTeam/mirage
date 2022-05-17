package optimizer

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"

	"github.com/chai2010/webp"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	log "github.com/sirupsen/logrus"
	giftowebp "github.com/sizeofint/gif-to-webp"
	"golang.org/x/image/bmp"
)

type Optimizer struct {
}

func NewOptimizer() *Optimizer {
	return &Optimizer{}
}

func (o *Optimizer) Optimize(data []byte, quality int64) (optimized []byte, originalContentType string, err error) {
	var buf bytes.Buffer
	contentType := http.DetectContentType(data)
	if strings.Contains(contentType, "gif") {
		converter := giftowebp.NewConverter()
		converter.LoopCompatibility = false
		converter.WebPConfig.SetQuality(float32(quality))
		converter.WebPConfig.SetMethod(4)
		webpBin, err := converter.Convert(data)
		if err != nil {
			return nil, contentType, errors.Err(err)
		}
		return webpBin, contentType, nil
	} else {
		img, err := readRawImage(data, contentType, 16383*16383)
		if err != nil {
			return nil, contentType, err
		}
		err = webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: float32(quality)})
		if err != nil {
			return nil, contentType, errors.Err(err)
		}
	}

	return buf.Bytes(), contentType, nil
}
func readRawImage(data []byte, contentType string, maxPixel int) (img image.Image, err error) {
	if strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg") {
		img, err = jpeg.Decode(bytes.NewReader(data))
	} else if strings.Contains(contentType, "png") {
		img, err = png.Decode(bytes.NewReader(data))
	} else if strings.Contains(contentType, "bmp") {
		img, err = bmp.Decode(bytes.NewReader(data))
	} else {
		return nil, errors.Err("%s type is not supported", contentType)
	}
	if err != nil || img == nil {
		errInfo := fmt.Sprintf("image file is corrupted: %v", err)
		log.Errorln(errInfo)
		return nil, errors.Err(errInfo)
	}

	x, y := img.Bounds().Max.X, img.Bounds().Max.Y
	if x > maxPixel || y > maxPixel {
		errInfo := fmt.Sprintf("WebP: (%dx%d) is too large", x, y)
		log.Warnf(errInfo)
		return nil, errors.Err(errInfo)
	}

	return img, nil
}