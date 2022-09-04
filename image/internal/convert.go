package internal

import (
	"image"
	"image/draw"
)

// ToNRGBA converts an image.Image to an [*image.NRGBA].
func ToNRGBA(img image.Image) *image.NRGBA {
	out := image.NewNRGBA(img.Bounds())
	draw.Draw(out, out.Bounds(), img, image.Point{}, draw.Src)
	return out
}
