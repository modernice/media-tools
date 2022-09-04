package image

import "image"

// Compressor compresses images.
type Compressor struct {
	compressor CompressionFunc
}

// CompressionFunc is the actual implementation of a compressor.
// Available implementations:
//   - [github.com/modernice/media-tools/image/compressor.JPEG]
type CompressionFunc func(image.Image) (*image.NRGBA, error)

// Compress returns a [*Compressor] that compresses images using the provided
// [CompressionFunc].
func Compress(compressor CompressionFunc) *Compressor {
	return &Compressor{compressor: compressor}
}

// Compress compresses an image using the configured [CompressionFunc].
func (c *Compressor) Compress(img image.Image) (*image.NRGBA, error) {
	return c.compressor(img)
}
