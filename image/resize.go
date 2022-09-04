package image

import (
	"image"

	"github.com/disintegration/imaging"
)

// Resizer resizes images to a set of dimensions.
type Resizer struct {
	dimensions []Dimensions
	filter     imaging.ResampleFilter
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

// DimensionProvider provides image dimensions to Resizer. DimensionProvider is
// implemented by [DimensionList] and [DimensionMap].
type DimensionProvider interface {
	Dimensions() []Dimensions
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

// Resize resizes an image to the configured dimensinos.
func (r *Resizer) Resize(img image.Image) ([]*image.NRGBA, error) {
	resized := make([]*image.NRGBA, len(r.dimensions))
	for i, dim := range r.dimensions {
		resized[i] = r.resize(img, dim)
	}
	return resized, nil
}

func (r *Resizer) resize(img image.Image, dim Dimensions) *image.NRGBA {
	return imaging.Resize(img, dim.Width(), dim.Height(), r.filter)
}

func (r *Resizer) Process(ctx ProcessorContext) ([]*image.NRGBA, error) {
	resized, err := r.Resize(ctx.Image())
	if err != nil {
		return nil, err
	}
	return append([]*image.NRGBA{ctx.Image()}, resized...), nil
}
