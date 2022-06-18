package helios

import (
	"errors"
	"fmt"
	i "image"
	"image/png"
	"os"
)

var ImageNotFoundError = errors.New("file not found")
var ImageUnsupportedFormatError = errors.New("unsupported format")
var ImageUnknownError = errors.New("unknown error")

type Image struct {
	img                 i.Image
	x                   float64
	y                   float64
	confidenceThreshold float64
}

func NewImage(path string, confidenceThreshold float64) (*Image, error) {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%w: %s", ImageNotFoundError, path)
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ImageUnknownError, err.Error())
	}

	defer func() { _ = file.Close() }()

	// Must specifically use jpeg.Decode() or it
	// would encounter unknown format error
	image, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ImageUnsupportedFormatError, err.Error())
	}

	return &Image{
		img:                 image,
		confidenceThreshold: confidenceThreshold,
	}, nil
}
