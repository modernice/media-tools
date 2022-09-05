package image

import (
	"fmt"
)

// Dimensions are the width and height of an image, in pixels.
type Dimensions [2]int

// Width returns the width of the image.
func (d Dimensions) Width() int {
	return d[0]
}

// Width returns the height of the image.
func (d Dimensions) Height() int {
	return d[1]
}

func (d Dimensions) String() string {
	return fmt.Sprintf("[width=%d, height=%d]", d.Width(), d.Height())
}

// DimensionList is a list of [Dimensions].
type DimensionList []Dimensions

// Dimensions returns the dimension list.
func (dl DimensionList) Dimensions() []Dimensions {
	return dl
}

// DimensionMap provides named [Dimensions].
type DimensionMap map[string]Dimensions

// Dimensions returns the dimension list.
func (dm DimensionMap) Dimensions() []Dimensions {
	dims := make([]Dimensions, 0, len(dm))
	for _, dim := range dm {
		dims = append(dims, dim)
	}
	return dims
}

// Tag returns the configured tag for the given [Dimensions].
func (dm DimensionMap) Tag(dim Dimensions) string {
	for tag, d := range dm {
		if d == dim {
			return tag
		}
	}
	return ""
}
