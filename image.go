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
	path                string
	confidenceThreshold float64
}

func (i *Image) GetImage() i.Image {
	return i.img
}

func (i *Image) GetPath() string {
	return i.path
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
		path:                path,
		confidenceThreshold: confidenceThreshold,
	}, nil
}
