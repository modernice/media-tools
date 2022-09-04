package internal

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"

	"github.com/vitali-fedulov/images4"
)

// ToNRGBA converts an image.Image to an [*image.NRGBA].
func ToNRGBA(img image.Image) *image.NRGBA {
	out := image.NewNRGBA(img.Bounds())
	draw.Draw(out, out.Bounds(), img, image.Point{}, draw.Src)
	return out
}

// SameImages returns whether two images are identical. The images are allowed
// to have different bounds and still be the same.
func SameImages(a, b image.Image) bool {
	iconA := images4.Icon(a)
	iconB := images4.Icon(b)
	return images4.Similar(iconA, iconB)
}

// EqualImages returns whether two images are identical, including theirs bounds.
func EqualImages(a, b image.Image) bool {
	return a.Bounds() == b.Bounds() && SameImages(a, b)
}

// SizeOf returns the JPEG size of an image in bytes.
func SizeOf(img image.Image) (int, error) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100}); err != nil {
		return 0, fmt.Errorf("encode as JPEG: %w", err)
	}
	return buf.Len(), nil
}
