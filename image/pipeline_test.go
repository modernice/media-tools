package image_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/modernice/media-tools/image"
	"github.com/modernice/media-tools/image/compressor"
	"github.com/modernice/media-tools/image/internal"
)

func TestPipeline_Run(t *testing.T) {
	pipe := image.Pipeline{
		image.Resize(image.DimensionMap{"sm": {360}, "md": {640}, "lg": {960}}),
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

	if !internal.SameImages(original, result.Images[0].Image) {
		t.Fatalf("first image should be the original\n%s", cmp.Diff(original, result.Images[0].Image))
	}

	if len(result.Images[0].Tags) != 1 {
		t.Fatalf("original image should have exactly 1 tag; has %d", len(result.Images[0].Tags))
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

	if !result.Images[1].Tags.Contains("compressor=jpeg,quality=75") {
		t.Fatalf("second image should have tag %q", "compressor=jpeg,quality=75")
	}

	if !result.Images[3].Tags.Contains("resized") {
		t.Fatalf("last image should have tag %q", "resized")
	}

	if !result.Images[3].Tags.Contains("size=lg") {
		t.Fatalf("last image should have tag %q", "size=lg")
	}

	if !result.Images[3].Tags.Contains("compressed") {
		t.Fatalf("last image should have tag %q", "compressed")
	}

	if !result.Images[3].Tags.Contains("compressor=jpeg,quality=75") {
		t.Fatalf("last image should have tag %q", "compressor=jpeg,quality=75")
	}
}

func TestPipeline_Run_CompressOriginal(t *testing.T) {
	pipe := image.Pipeline{
		image.Resize(image.DimensionMap{"sm": {360}, "md": {640}, "lg": {960}}),
		image.Compress(compressor.JPEG(75), image.CompressOriginal(true)),
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

	if !result.Images[0].Tags.Contains("compressor=jpeg,quality=75") {
		t.Fatalf("original image should have tag %q", "compressor=jpeg,quality=75")
	}
}

func TestPipeline_Run_ResizeDiscardInput(t *testing.T) {
	pipe := image.Pipeline{
		image.Resize(image.DimensionMap{"sm": {360}, "md": {640}, "lg": {960}}, image.DiscardInput(true)),
		image.Compress(compressor.JPEG(75)),
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
