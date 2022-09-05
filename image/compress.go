package image

import (
	"image"

	"github.com/modernice/media-tools/image/internal"
)

// Compressor compresses images.
type Compressor struct {
	compressor       CompressionFunc
	compressOriginal bool
}

// CompressionFunc is the actual implementation of a compressor.
// Available implementations:
//   - [github.com/modernice/media-tools/image/compressor.JPEG]
type CompressionFunc func(image.Image) (image.Image, error)

// CompressorOption is an option for a [*Compressor].
type CompressorOption func(*Compressor)

// CompressOriginal returns a CompressorOption that enables compression of
// the original image of a [Pipeline]. By default, the original image will not
// be compressed and instead returned as is to preserve the original image
// quality.
func CompressOriginal(v bool) CompressorOption {
	return func(c *Compressor) {
		c.compressOriginal = v
	}
}

// Compress returns a [*Compressor] that compresses images using the provided
// [CompressionFunc].
func Compress(compressor CompressionFunc, opts ...CompressorOption) *Compressor {
	c := &Compressor{compressor: compressor}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Compress compresses an image using the configured [CompressionFunc].
func (c *Compressor) Compress(img image.Image) (image.Image, error) {
	compressed, err := c.compressor(img)
	if err != nil {
		return nil, err
	}
	return internal.ToNRGBA(compressed), nil
}

// Process implements [Processor]. By default, the original image will not be
// compressed and returned as is to preserve quality. To also compress the
// original image, pass the [CompressOriginal] option to [Compress].
func (c *Compressor) Process(ctx ProcessorContext) ([]image.Image, error) {
	if !c.compressOriginal && ctx.Original() {
		return []image.Image{ctx.Image()}, nil
	}

	compressed, err := c.Compress(ctx.Image())
	if err != nil {
		return nil, err
	}

	return []image.Image{compressed}, nil
}
