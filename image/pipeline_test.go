package image_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/modernice/media-tools/image"
	"github.com/modernice/media-tools/image/compression"
	"github.com/modernice/media-tools/image/internal"
)

func TestPipeline_Run(t *testing.T) {
	pipe := image.Pipeline{
		image.Resize(image.DimensionMap{"sm": {360}, "md": {640}, "lg": {960}}),
		image.CompressMany([]image.Compression{
			compression.JPEG(75),
			compression.JPEG(50),
		}),
		image.Tag(image.NewTags("foo", "bar")),
		image.TagBy(func(p image.Processed) image.Tags {
			suffix := "non-original"
			if p.Original {
				suffix = "original"
			}
			return image.NewTags(fmt.Sprintf("tagby:%s", suffix))
		}),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	original := newExample()

	result, err := pipe.Run(ctx, original)
	if err != nil {
		t.Fatalf("run pipeline: %v", err)
	}

	if len(result.Images) != 7 {
		t.Fatalf("pipeline should return 7 images (including the original); got %d", len(result.Images))
	}

	if !internal.SameImages(original, result.Images[0].Image) {
		t.Fatalf("first image should be the original\n%s", cmp.Diff(original, result.Images[0].Image))
	}

	if len(result.Images[0].Tags) != 4 {
		t.Fatalf("original image should have exactly 4 tags; has %d", len(result.Images[0].Tags))
	}

	if !result.Images[0].Original {
		t.Fatalf("first image should be the original")
	}

	if !result.Images[0].Tags.Contains(image.Original) {
		t.Fatalf("original image should have tag %q", image.Original)
	}

	if !result.Images[1].Tags.Contains("resized") {
		t.Fatalf("second image should have tag %q", "resized")
	}

	if !result.Images[1].Tags.Contains("size=sm") {
		t.Fatalf("second image should have tag %q", "size=sm")
	}

	if !result.Images[1].Tags.Contains("compressed") {
		t.Fatalf("second image should have tag %q", "compressed")
	}

	if !result.Images[1].Tags.Contains("compression=jpeg,quality=75") {
		t.Fatalf("second image should have tag %q", "compression=jpeg,quality=75")
	}

	lastIdx := len(result.Images) - 1
	secondLastIdx := len(result.Images) - 2

	if !result.Images[secondLastIdx].Tags.Contains("resized") {
		t.Fatalf("second-to-last image should have tag %q", "resized")
	}

	if !result.Images[secondLastIdx].Tags.Contains("size=lg") {
		t.Fatalf("second-to-last image should have tag %q", "size=lg")
	}

	if !result.Images[secondLastIdx].Tags.Contains("compressed") {
		t.Fatalf("second-to-last image should have tag %q", "compressed")
	}

	if !result.Images[secondLastIdx].Tags.Contains("compression=jpeg,quality=75") {
		t.Fatalf("second-to-last image should have tag %q", "compression=jpeg,quality=75")
	}

	if !result.Images[lastIdx].Tags.Contains("resized") {
		t.Fatalf("last image should have tag %q", "resized")
	}

	if !result.Images[lastIdx].Tags.Contains("size=lg") {
		t.Fatalf("last image should have tag %q", "size=lg")
	}

	if !result.Images[lastIdx].Tags.Contains("compressed") {
		t.Fatalf("last image should have tag %q", "compressed")
	}

	if !result.Images[lastIdx].Tags.Contains("compression=jpeg,quality=50") {
		t.Fatalf("last image should have tag %q", "compression=jpeg,quality=50")
	}

	for _, img := range result.Images {
		if !img.Tags.Contains("foo") {
			t.Fatalf("all images should have tag %q", "foo")
		}

		if !img.Tags.Contains("bar") {
			t.Fatalf("all images should have tag %q", "bar")
		}

		if img.Original && !img.Tags.Contains("tagby:original") {
			t.Fatalf("original image should have tag %q", "tagby:original")
		}

		if !img.Original && !img.Tags.Contains("tagby:non-original") {
			t.Fatalf("non-original image should have tag %q", "tagby:non-original")
		}
	}
}

func TestPipeline_Run_CompressOriginal(t *testing.T) {
	pipe := image.Pipeline{
		image.Resize(image.DimensionMap{"sm": {360}, "md": {640}, "lg": {960}}),
		image.Compress(compression.JPEG(75), image.CompressOriginal(true)),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	original := newExample()

	result, err := pipe.Run(ctx, original)
	if err != nil {
		t.Fatalf("run pipeline: %v", err)
	}

	if !result.Images[0].Original {
		t.Fatalf("first image should be the original")
	}

	if !result.Images[0].Tags.Contains("compressed") {
		t.Fatalf("original image should have tag %q", "compressed")
	}

	if !result.Images[0].Tags.Contains("compression=jpeg,quality=75") {
		t.Fatalf("original image should have tag %q", "compression=jpeg,quality=75")
	}
}

func TestPipeline_Run_ResizeDiscardInput(t *testing.T) {
	pipe := image.Pipeline{
		image.Resize(image.DimensionMap{"sm": {360}, "md": {640}, "lg": {960}}, image.DiscardInput(true)),
		image.Compress(compression.JPEG(75)),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	original := newExample()

	result, err := pipe.Run(ctx, original)
	if err != nil {
		t.Fatalf("run pipeline: %v", err)
	}

	if len(result.Images) != 3 {
		t.Fatalf("pipeline should return 3 images (excluding the original); got %d", len(result.Images))
	}

	if result.Images[0].Original {
		t.Fatalf("original image should be discarded")
	}

	if !result.Images[0].Tags.Contains("resized") {
		t.Fatalf("first image should have tag %q", "resized")
	}

	if !result.Images[0].Tags.Contains("size=sm") {
		t.Fatalf("first image should have tag %q", "size=sm")
	}

	if !result.Images[0].Tags.Contains("compressed") {
		t.Fatalf("first image should have tag %q", "compressed")
	}
}
