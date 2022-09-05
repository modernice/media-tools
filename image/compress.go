package image

import (
	"image"

	"github.com/modernice/media-tools/image/internal"
)

var _ Processor = (*Compressor)(nil)

const (
	// Compressed is the tag that is assigned to compressed images.
	Compressed = "compressed"

	// AnonymousCompression is the name [Compression.Name] of a CompressionFunc.
	AnonymousCompression = "anonymous"
)

// Compressor compresses images.
type Compressor struct {
	compressions     []Compression
	compressOriginal bool
}

// Compression provides the actual implementation for compressing images.
type Compression interface {
	// Compress compresses an image.
	Compress(img image.Image) (image.Image, error)
}

// CompressionFunc allow a function to be used as a [Compression].
type CompressionFunc func(image.Image) (image.Image, error)

func (fn CompressionFunc) Compress(img image.Image) (image.Image, error) {
	return fn(img)
}

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
// [Compression]. If the provided [Compression] has a `Tags() Tags` method,
// the returned [*Compressor] will append these tags to compressed images
// when calling [*Compressor.Process].
func Compress(compression Compression, opts ...CompressorOption) *Compressor {
	return CompressMany([]Compression{compression}, opts...)
}

// CompressMany returns a [*Compressor] that compresses images using the provided
// [Compression]s. If a provided [Compression] method has a `Tags() Tags` method,
// the returned [*Compressor] will append these tags to compressed image when
// calling [*Compressor.Process].
func CompressMany(compressions []Compression, opts ...CompressorOption) *Compressor {
	c := &Compressor{compressions: compressions}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Compress compresses an image using the configured [Compression]s.
func (c *Compressor) Compress(img image.Image) ([]image.Image, error) {
	out := make([]image.Image, len(c.compressions))
	for i, compression := range c.compressions {
		compressed, err := compression.Compress(img)
		if err != nil {
			return nil, err
		}
		out[i] = internal.ToNRGBA(compressed)
	}
	return out, nil
}

// Process implements [Processor]. By default, the original image will not be
// compressed and returned as is to preserve quality. To also compress the
// original image, pass the [CompressOriginal] option to [Compress].
func (c *Compressor) Process(ctx ProcessorContext) ([]Processed, error) {
	pimg := ctx.Image()

	if !c.compressOriginal && pimg.Original {
		return []Processed{pimg}, nil
	}

	compressed, err := c.Compress(ctx.Image().Image)
	if err != nil {
		return nil, err
	}

	out := make([]Processed, len(c.compressions))
	for i, compression := range c.compressions {
		var compressionTags Tags
		if tagger, isTagger := compression.(interface{ Tags() Tags }); isTagger {
			compressionTags = tagger.Tags()
		}

		out[i] = Processed{
			Image:    compressed[i],
			Tags:     pimg.Tags.With(Compressed).With(compressionTags...),
			Original: pimg.Original,
		}
	}

	return out, nil
}
