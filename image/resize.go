package image

import (
	"fmt"
	"image"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/modernice/media-tools/internal/slices"
	stdslices "golang.org/x/exp/slices"
)

var _ Processor = (*Resizer)(nil)

// Resized is the tag that is assigned to resized images.
const Resized = "resized"

// Resizer resizes images to a set of dimensions.
type Resizer struct {
	dimensionProvider DimensionProvider
	dimensions        DimensionList
	filter            imaging.ResampleFilter
	discardInput      bool
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

// DiscardInput returns a ResizerOption that discards the input image from the
// resize result when executed in a [Pipeline].
func DiscardInput(v bool) ResizerOption {
	return func(r *Resizer) {
		r.discardInput = v
	}
}

// Resize returns a Resizer that resizes images to the given dimensions.
func Resize(dimensions DimensionProvider, opts ...ResizerOption) *Resizer {
	r := &Resizer{
		dimensionProvider: dimensions,
		filter:            imaging.Lanczos,
	}

	for _, opt := range opts {
		opt(r)
	}

	r.dimensions = r.dimensionProvider.Dimensions()
	stdslices.SortFunc(r.dimensions, func(a, b Dimensions) int {
		if a.Width() < b.Width() {
			return -1
		}

		if a.Width() > b.Width() {
			return 1
		}

		if a.Height() < b.Height() {
			return -1
		}

		if a.Height() > b.Height() {
			return 1
		}

		return 0
	})

	return r
}

// Resize resizes an image to the configured dimensinos. The input image is not
// returned in the result.
func (r *Resizer) Resize(img image.Image) ([]image.Image, error) {
	resized, err := r.resizeInternal(img)
	if err != nil {
		return nil, err
	}
	return slices.Map(func(r resizedImage) image.Image { return r.image }, resized), nil
}

type resizedImage struct {
	image      image.Image
	dimensions Dimensions
}

func (r *Resizer) resizeInternal(img image.Image) ([]resizedImage, error) {
	resized := make([]resizedImage, len(r.dimensions))
	for i, dim := range r.dimensions {
		resized[i] = resizedImage{
			image:      r.resize(img, dim),
			dimensions: dim,
		}
	}
	return resized, nil
}

func (r *Resizer) resize(img image.Image, dim Dimensions) image.Image {
	return imaging.Resize(img, dim.Width(), dim.Height(), r.filter)
}

// Process implements [Processor]. The input image is returned in the result as
// the first element.
func (r *Resizer) Process(ctx ProcessorContext) ([]Processed, error) {
	input := ctx.Image()

	resized, err := r.resizeInternal(input.Image)
	if err != nil {
		return nil, err
	}

	tagger, isTagger := r.dimensionProvider.(interface{ Tag(Dimensions) string })

	processed := make([]Processed, len(resized))
	baseTags := input.Tags.Without(Original)
	for i, rimg := range resized {
		tags := baseTags.With("resized")

		if isTagger {
			tags = tags.With(fmt.Sprintf("size=%s", tagger.Tag(rimg.dimensions)))
		}

		processed[i] = Processed{
			Image: rimg.image,
			Tags:  tags,
		}
	}

	if r.discardInput {
		return processed, nil
	}

	return append([]Processed{ctx.Image()}, processed...), nil
}

// DimensionName extracts the dimension name from the tags of a processed image.
func DimensionName(tags Tags) string {
	for _, tag := range tags {
		if strings.HasPrefix(tag, "size=") {
			return tag[5:]
		}
	}
	return ""
}
