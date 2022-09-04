package image_test

import (
	"bytes"
	"fmt"
	stdimage "image"
	"image/jpeg"
	"testing"

	"github.com/modernice/media-tools/image"
	"github.com/modernice/media-tools/image/compressor"
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

			bounds := img.Bounds()
			saveCompressed(t, quality, image.Dimensions{bounds.Dx(), bounds.Dy()}, compressed)
		})
	}
}

func getImageSize(t *testing.T, img stdimage.Image) int {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100}); err != nil {
		t.Fatalf("encode as JPEG: %v", err)
	}
	return buf.Len()
}

func saveCompressed(t *testing.T, quality int, dim image.Dimensions, img *stdimage.NRGBA) {
	saveOutImage(t, fmt.Sprintf("compressed-%dx%d-%d.jpg", dim.Width(), dim.Height(), quality), dim, img)
}
