package util

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"golang.org/x/image/draw"
)

// Readable image formats: jpeg, png
func LoadImage(data []byte, fileType string) (image.Image, error) {
	var img image.Image
	var err error

	switch fileType {
	case "image/jpeg":
		img, err = jpeg.Decode(bytes.NewReader(data))
	case "image/png":
		img, err = png.Decode(bytes.NewReader(data))
	default:
		err = fmt.Errorf("unsupported file type")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to decode the image: %w", err)
	}

	return img, nil
}

// Aspect ratio preserving image resizing
func Resize(data []byte, maxHeight int) []byte {
	img, err := LoadImage(data, "image/jpeg")
	if err != nil {
		return nil
	}

	// Calculate the new size
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	if width > height {
		width = width * maxHeight / height
		height = maxHeight
	} else {
		height = height * maxHeight / width
		width = maxHeight
	}

	// Resize
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	// Encode the image
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, dst, nil)

	return buf.Bytes()
}

// func ConvertToAVIF(img image.Image) ([]byte, error) {
// 	quality := 30

// 	buf := new(bytes.Buffer)
// 	err := avif.Encode(buf, img, avif.Options{Quality: quality})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to encode the image to AVIF: %w", err)
// 	}

// 	return buf.Bytes(), nil
// }
