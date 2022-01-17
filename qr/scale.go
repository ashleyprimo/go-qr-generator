package qr

import (
	"errors"
	"image"

	"github.com/nfnt/resize"

	log "github.com/sirupsen/logrus"
)

func scale(code image.Image, sizeInt uint) (image.Image, error) {
	log.Debugf("Scaling Image")

	// Resize 'image.Image' to requested size
	img := resize.Resize(sizeInt, 0, code, resize.NearestNeighbor)
	if img != nil {
		return img, nil
	} else {
		return nil, errors.New("Failed to resize image")
	}
}
