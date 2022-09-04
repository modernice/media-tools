package image_test

import (
	"context"
	stdimage "image"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/modernice/media-tools/image"
	"github.com/modernice/media-tools/image/compressor"
	"github.com/vitali-fedulov/images4"
)

func TestPipeline(t *testing.T) {
	pipe := image.Pipeline{
		image.Resize(image.DimensionList{{360}, {640}, {960}}),
		image.Compress(compressor.JPEG(75)),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	original := newExample()

	result, err := pipe.Run(ctx, original)
	if err != nil {
		t.Fatalf("run pipeline: %v", err)
	}

	if len(result.Images) != 4 {
		t.Fatalf("pipeline should return 4 images (including the original); got %d", len(result.Images))
	}

	for _, img := range result.Images {
		if img.ID == "" {
			t.Fatal("processed images should have an id")
		}
	}

	if !equalImages(original, result.Images[0].Image) {
		t.Fatalf("first image should be the original\n%s", cmp.Diff(original, result.Images[0].Image))
	}
}

func equalImages(a, b stdimage.Image) bool {
	iconA := images4.Icon(a)
	iconB := images4.Icon(b)
	return images4.Similar(iconA, iconB)
}
