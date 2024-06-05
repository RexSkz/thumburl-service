package imageservice

import (
	"bytes"
	"image"

	"github.com/chai2010/webp"
	"golang.org/x/image/draw"
)

func ResizeWebp(webpData []byte, width, height int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(webpData))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	wRatio := float64(bounds.Dx()) / float64(width)
	hRatio := float64(bounds.Dy()) / float64(height)
	ratio := wRatio
	if hRatio > wRatio {
		ratio = hRatio
	}

	newWidth := int(float64(bounds.Dx()) / ratio)
	newHeight := int(float64(bounds.Dy()) / ratio)
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(newImg, newImg.Bounds(), img, bounds, draw.Over, nil)

	result := new(bytes.Buffer)
	if err = webp.Encode(result, newImg, &webp.Options{Quality: 90}); err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}
