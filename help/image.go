package help

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
)

// ToPng converts an image to png
func ToPng(imageBytes []byte) ([]byte, error) {
	contentType := http.DetectContentType(imageBytes)

	switch contentType {
	case "image/png":
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, err
		}

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}
	return nil, fmt.Errorf("unable to convert %#v to png", contentType)
}
