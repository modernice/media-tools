package image_test

import (
	"context"
	"testing"

	"github.com/modernice/media-tools/image"
	"github.com/modernice/media-tools/image/internal"
)

func TestTagger_Process(t *testing.T) {
	original := newExample()

	tagger := image.Tag(image.NewTags("foo", "bar"))

	ctx := image.NewProcessorContext(context.Background(), image.Processed{Image: original, Original: true})

	tagged, err := tagger.Process(ctx)
	if err != nil {
		t.Fatalf("run processor: %v", err)
	}

	if len(tagged) != 1 {
		t.Fatalf("expected 1 image; got %d", len(tagged))
	}

	if !internal.EqualImages(original, tagged[0].Image) {
		t.Fatalf("tagger should not alter the image")
	}

	if !tagged[0].Tags.Contains("foo") {
		t.Fatalf("expected image to have tag %q", "foo")
	}

	if !tagged[0].Tags.Contains("bar") {
		t.Fatalf("expected image to have tag %q", "bar")
	}
}
