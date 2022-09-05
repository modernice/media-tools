package compression

import (
	"bytes"
	"fmt"
	stdimage "image"
	"image/jpeg"

	"github.com/modernice/media-tools/image"
	"github.com/modernice/media-tools/image/internal"
)

// JPEG retrurns an [image.Compression] that compresses images using the JPEG
// encoder's "quality" option.
func JPEG(quality int) image.Compression {
	return &jpegCompression{quality: quality}
}

type jpegCompression struct{ quality int }

func (jc *jpegCompression) Compress(img stdimage.Image) (stdimage.Image, error) {
	var buf bytes.Buffer

	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: jc.quality}); err != nil {
		return nil, fmt.Errorf("encode as JPEG: %w", err)
	}

	decoded, err := jpeg.Decode(&buf)
	if err != nil {
		return nil, fmt.Errorf("decode JPEG: %w", err)
	}

	return internal.ToNRGBA(decoded), nil
}

// Tags returns the tags that should be assigned to images that are compressed
// by the JPEG compression.
func (jc *jpegCompression) Tags() image.Tags {
	return image.NewTags(fmt.Sprintf("compressor=jpeg,quality=%d", jc.quality))
}
