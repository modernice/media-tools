# Image Processing Pipeline

Package `image` provides a Pluggable Image Processor. Multiple Processors can be
chained together to form a Pipeline, which can be executed to arbitrarily
process images. This can be used to, for example, resize and compress images in
a web application that supports image uploads.

## Usage

```go
package example

import (
	stdimage "image"
	"github.com/modernice/media-tools/image"
)

func example(img stdimage.Image) {
	pipeline := image.Pipeline{
		// First, resize the image to 5 named dimensions.
		image.Resize(image.DimensionMap{
			"xs": {360},
			"sm": {640},
			"md": {960},
			"lg": {1280},
			"xl": {1920},
		}),

		// Then, compress the original and each resized image
		// to two qualities.
		image.Compress(compression.JPEG(80), compression.JPEG(50)),
	}

	result, err := pipeline.Run(context.TODO(), img)
	if err != nil {
		panic(err)
	}

	// 1 original + 1 original compressed
	// + 5 resized + 5 resized-compressed
	len(result.Images) == 12

	compressed := result.Find("compressed")
	len(compressed) == 6

	image.CompressionName(compressed[0].Tags) == "jpeg"
	image.CompressionQuality(compressed[0].Tags) == 80

	resized := result.Find("resized")
	len(resized) == 6

	image.DimensionName(resized[0].Tags) == "xs"
	image.DimensionName(resized[1].Tags) == "sm"
	image.DimensionName(resized[4].Tags) == "xl"

	sm := result.Find("size=lg")
	len(sm) == 1

	image.DimensionName(sm[0].Tags) == "sm"
}
```
