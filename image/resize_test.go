package image_test

import (
	_ "embed"
	"fmt"
	stdimage "image"
	"testing"

	"github.com/modernice/media-tools/image"
)

func TestResizer_Resize_DimensionList(t *testing.T) {
	dimensions := image.DimensionList{
		{100, 100},
		{200, 200},
		{300, 500},
		{960},
		{1280},
		{1920},
	}

	resizer := image.Resize(dimensions)

	img := newExample()

	resized, err := resizer.Resize(img)
	if err != nil {
		t.Fatalf("failed to resize: %v", err)
	}

	if len(resized) != len(dimensions) {
		t.Fatalf("expected %d resized images; got %d", len(dimensions), len(resized))
	}

	for i, dim := range dimensions {
		rimg := resized[i]

		if rimg.Bounds().Dx() != dim.Width() {
			t.Fatalf("expected width %d; got %d", dim.Width(), rimg.Bounds().Dx())
		}

		if dim.Height() > 0 && rimg.Bounds().Dy() != dim.Height() {
			t.Fatalf("expected height %d; got %d", dim.Height(), rimg.Bounds().Dy())
		}

		saveResized(t, dim, rimg)
	}
}

func saveResized(t *testing.T, dim image.Dimensions, img *stdimage.NRGBA) {
	saveOutImage(t, fmt.Sprintf("resized-%dx%d.jpg", dim.Width(), dim.Height()), dim, img)
}
