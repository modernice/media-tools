package image_test

import (
	"bytes"
	_ "embed"
	"fmt"
	stdimage "image"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/modernice/media-tools/image"
)

//go:embed testdata/example.jpg
var example []byte

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

	outputDir := prepareOutputDir(t)

	for i, dim := range dimensions {
		rimg := resized[i]

		if rimg.Bounds().Dx() != dim.Width() {
			t.Fatalf("expected width %d; got %d", dim.Width(), rimg.Bounds().Dx())
		}

		if dim.Height() > 0 && rimg.Bounds().Dy() != dim.Height() {
			t.Fatalf("expected height %d; got %d", dim.Height(), rimg.Bounds().Dy())
		}

		if err := saveResized(t, outputDir, dim, rimg); err != nil {
			t.Fatalf("save resizd image: %v", err)
		}
	}
}

func newExample() stdimage.Image {
	img, err := jpeg.Decode(bytes.NewReader(example))
	if err != nil {
		panic(err)
	}
	return img
}

func prepareOutputDir(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	outputDir := filepath.Join(wd, "testoutput")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("create %s: %v", outputDir, err)
	}

	return outputDir
}

func saveResized(t *testing.T, outDir string, dim image.Dimensions, img *stdimage.NRGBA) error {
	path := filepath.Join(outDir, fmt.Sprintf("resized-%d-%d.jpg", dim.Width(), dim.Height()))
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create %s: %v", path, err)
	}
	defer f.Close()

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100}); err != nil {
		t.Fatalf("encode resized image: %v", err)
	}

	if _, err := f.Write(buf.Bytes()); err != nil {
		t.Fatalf("write resized image: %v", err)
	}

	if err := f.Close(); err != nil {
		t.Fatalf("close %s: %v", path, err)
	}

	return nil
}
