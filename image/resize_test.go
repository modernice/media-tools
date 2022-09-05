package image_test

import (
	"context"
	_ "embed"
	"fmt"
	stdimage "image"
	"testing"

	"github.com/modernice/media-tools/image"
	"github.com/modernice/media-tools/image/internal"
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

func TestResizer_Process(t *testing.T) {
	dimensions := image.DimensionList{{360}, {640}, {960}}

	resizer := image.Resize(dimensions)

	original := newExample()
	ctx := image.NewProcessorContext(context.Background(), image.Processed{Image: original, Original: true})

	resized, err := resizer.Process(ctx)
	if err != nil {
		t.Fatalf("run processor: %v", err)
	}

	if len(resized) != len(dimensions)+1 {
		t.Fatalf("expected %d resized images (including the original); got %d", len(dimensions)+1, len(resized))
	}

	if !internal.EqualImages(original, resized[0].Image) {
		t.Fatalf("first returned image should be the original")
	}
}

func TestResizer_Process_DiscardOriginal(t *testing.T) {
	dimensions := image.DimensionList{{360}, {640}, {960}}

	resizer := image.Resize(dimensions, image.DiscardInput(true))

	original := newExample()
	ctx := image.NewProcessorContext(context.Background(), image.Processed{Image: original, Original: true})

	resized, err := resizer.Process(ctx)
	if err != nil {
		t.Fatalf("run processor: %v", err)
	}

	if len(resized) != len(dimensions) {
		t.Fatalf("expected %d resized images (excluding the original); got %d", len(dimensions), len(resized))
	}

	if internal.EqualImages(original, resized[0].Image) {
		t.Fatalf("original image should not be returned")
	}
}

// When running an [*image.Resizer] in an [image.Pipeline] with an
// [image.DimensionMap], the [image.PipelineResult] must provide the correct
// tags for the processed images.
func TestResizer_Process_Pipeline_DimensionMap(t *testing.T) {
	dimensions := image.DimensionMap{
		"thumb": {360},
		"sm":    {640},
		"md":    {960},
		"lg":    {1280},
		"xl":    {1920},
	}

	original := newExample()

	pipe := image.Pipeline{image.Resize(dimensions)}

	result, err := pipe.Run(context.Background(), original)
	if err != nil {
		t.Fatalf("run pipeline: %v", err)
	}

	orgTags := result.Images[0].Tags
	if len(orgTags) != 1 {
		t.Fatalf("original image should have exactly 1 tag")
	}

	if !orgTags.Contains(image.Original) {
		t.Fatalf("original image should have tag %q", image.Original)
	}

	thumbTags := result.Images[1].Tags

	if thumbTags.Contains(image.Original) {
		t.Fatalf("resized image must not have the %q tag", image.Original)
	}

	if !thumbTags.Contains("resized") {
		t.Fatalf("thumb image must have tag %q", "resized")
	}

	if !thumbTags.Contains("size=thumb") {
		t.Fatalf("thumb image must have tag %q", "size=thumb")
	}

	if name := image.DimensionName(result.Images[0].Tags); name != "" {
		t.Fatalf("original image must not have a dimension name; got %q", name)
	}

	if name := image.DimensionName(result.Images[1].Tags); name != "thumb" {
		t.Fatalf("thumb image must not dimension name %q; got %q", "thumb", name)
	}

	if name := image.DimensionName(result.Images[5].Tags); name != "xl" {
		t.Fatalf("xl image must not have dimension name %q; got %q", "xl", name)
	}
}

func saveResized(t *testing.T, dim image.Dimensions, img stdimage.Image) {
	saveOutImage(t, fmt.Sprintf("resized-%dx%d.jpg", dim.Width(), dim.Height()), img)
}
