package main

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/gen2brain/avif"
	_ "github.com/gen2brain/heic"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

// Readable image formats
// jpeg, png, webp, heic
func loadImage(r io.Reader) (image.Image, string, error) {
	img, format, err := image.Decode(r)
	if err != nil {
		return nil, "", err
	}

	return img, format, nil
}

// Aspect ratio preserving image resizing
func resize(img image.Image, maxHeight int) image.Image {
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

	return dst
}

func convertToAVIF(img image.Image) (*bytes.Buffer, error) {
	quality := 30

	buf := new(bytes.Buffer)
	err := avif.Encode(buf, img, avif.Options{Quality: quality})
	if err != nil {
		return nil, err
	}

	return buf, nil
}
