package image

import (
	"image"

	"github.com/disintegration/imaging"
)

// Resizer resizes images to a set of dimensions.
type Resizer struct {
	dimensions      []Dimensions
	filter          imaging.ResampleFilter
	discardOriginal bool
}

// DimensionProvider provides image dimensions to Resizer. DimensionProvider is
// implemented by [DimensionList] and [DimensionMap].
type DimensionProvider interface {
	Dimensions() []Dimensions
}

// ResizerOption is an option for a Resizer.
type ResizerOption func(*Resizer)

// ResampleFilter returns a ResizerOption that sets the [imaging.ResampleFilter]
// to use when resizing images. Defaults to [imaging.Lanczos].
func ResampleFilter(filter imaging.ResampleFilter) ResizerOption {
	return func(r *Resizer) {
		r.filter = filter
	}
}

// DiscardOriginal returns a ResizerOption that discards the original image when
// executed in a [Pipeline].
func DiscardOriginal(v bool) ResizerOption {
	return func(r *Resizer) {
		r.discardOriginal = v
	}
}

// Resize returns a Resizer that resizes images to the given dimensions.
func Resize(dimensions DimensionProvider, opts ...ResizerOption) *Resizer {
	dims := dimensions.Dimensions()
	r := &Resizer{dimensions: dims, filter: imaging.Lanczos}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Resize resizes an image to the configured dimensinos. The input image is not
// returned in the result.
func (r *Resizer) Resize(img image.Image) ([]image.Image, error) {
	resized := make([]image.Image, len(r.dimensions))
	for i, dim := range r.dimensions {
		resized[i] = r.resize(img, dim)
	}
	return resized, nil
}

func (r *Resizer) resize(img image.Image, dim Dimensions) image.Image {
	return imaging.Resize(img, dim.Width(), dim.Height(), r.filter)
}

// Process implements [Processor]. The input image is returned in the result as
// the first element.
func (r *Resizer) Process(ctx ProcessorContext) ([]image.Image, error) {
	resized, err := r.Resize(ctx.Image())
	if err != nil {
		return nil, err
	}

	if r.discardOriginal {
		return resized, nil
	}

	return append([]image.Image{ctx.Image()}, resized...), nil
}
