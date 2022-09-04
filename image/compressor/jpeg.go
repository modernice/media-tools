package compressor

import (
	"bytes"
	"fmt"
	stdimage "image"
	"image/jpeg"

	"github.com/modernice/media-tools/image"
	"github.com/modernice/media-tools/image/internal"
)

// JPEG returns an [image.CompressionFunc] that compresses images using the JPEG
// encoder's "quality" option.
func JPEG(quality int) image.CompressionFunc {
	return func(img stdimage.Image) (*stdimage.NRGBA, error) {
		var buf bytes.Buffer

		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
			return nil, fmt.Errorf("encode JPEG: %w", err)
		}

		decoded, err := jpeg.Decode(&buf)
		if err != nil {
			return nil, fmt.Errorf("decode JPEG: %w", err)
		}

		return internal.ToNRGBA(decoded), nil
	}
}
