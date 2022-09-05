package image_test

import (
	"bytes"
	"context"
	"fmt"
	stdimage "image"
	"image/jpeg"
	"testing"

	"github.com/modernice/media-tools/image"
	"github.com/modernice/media-tools/image/compressor"
	"github.com/modernice/media-tools/image/internal"
)

func TestCompressor_Compress(t *testing.T) {
	qualities := []int{100, 50, 20, 5}

	for _, quality := range qualities {
		t.Run(fmt.Sprintf("quality=%d", quality), func(t *testing.T) {
			compressor := image.Compress(compressor.JPEG(quality))

			img := newExample()

			compressed, err := compressor.Compress(img)
			if err != nil {
				t.Fatalf("compress image: %v", err)
			}

			if compressed.Bounds() != img.Bounds() {
				t.Fatalf("compressed image has wrong bounds; want %v; got %v", img.Bounds(), compressed.Bounds())
			}

			orgSize := getImageSize(t, img)
			compressedSize := getImageSize(t, compressed)

			if quality < 100 && compressedSize >= orgSize {
				t.Fatalf("compressed image should be smaller than original; original=%d; compressed=%d", orgSize, compressedSize)
			}

			saveCompressed(t, quality, compressed)
		})
	}
}

func TestCompressor_Process_original(t *testing.T) {
	compressor := image.Compress(compressor.JPEG(50))

	original := newExample()
	ctx := image.NewProcessorContext(context.Background(), image.Processed{Image: original, Original: true})

	compressed, err := compressor.Process(ctx)
	if err != nil {
		t.Fatalf("run processor: %v", err)
	}

	if len(compressed) != 1 {
		t.Fatalf("expected 1 compressed image; got %d", len(compressed))
	}

	orgSize, err := internal.SizeOf(original)
	if err != nil {
		t.Fatalf("get original size: %v", err)
	}
	compressedSize, err := internal.SizeOf(compressed[0].Image)
	if err != nil {
		t.Fatalf("get compressed size: %v", err)
	}

	if orgSize != compressedSize {
		t.Fatalf("original image should not be compressed")
	}
}

func TestCompressor_Process_CompressOriginal(t *testing.T) {
	compressor := image.Compress(compressor.JPEG(50), image.CompressOriginal(true))

	original := newExample()
	ctx := image.NewProcessorContext(context.Background(), image.Processed{Image: original, Original: true})

	compressed, err := compressor.Process(ctx)
	if err != nil {
		t.Fatalf("run processor: %v", err)
	}

	if len(compressed) != 1 {
		t.Fatalf("expected 1 compressed image; got %d", len(compressed))
	}

	orgSize, err := internal.SizeOf(original)
	if err != nil {
		t.Fatalf("get original size: %v", err)
	}
	compressedSize, err := internal.SizeOf(compressed[0].Image)
	if err != nil {
		t.Fatalf("get compressed size: %v", err)
	}

	if orgSize == compressedSize {
		t.Fatalf("original image should be compressed when providing the CompressOriginal(true) option")
	}
}

func getImageSize(t *testing.T, img stdimage.Image) int {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100}); err != nil {
		t.Fatalf("encode as JPEG: %v", err)
	}
	return buf.Len()
}

func saveCompressed(t *testing.T, quality int, img stdimage.Image) {
	saveOutImage(t, fmt.Sprintf("compressed-%d.jpg", quality), img)
}
