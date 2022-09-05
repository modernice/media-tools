package image_test

import (
	"bytes"
	_ "embed"
	stdimage "image"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/modernice/media-tools/image/internal"
)

//go:embed testdata/example.jpg
var example []byte

var wd, _ = os.Getwd()

func newExample() *stdimage.NRGBA {
	img, err := jpeg.Decode(bytes.NewReader(example))
	if err != nil {
		panic(err)
	}
	return internal.ToNRGBA(img)
}

func saveOutImage(t *testing.T, name string, img stdimage.Image) {
	prepareOutputDir(t)

	path := filepath.Join(wd, "testoutput", name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create %s: %v", path, err)
	}
	defer f.Close()

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100}); err != nil {
		t.Fatalf("encode image: %v", err)
	}

	if _, err := f.Write(buf.Bytes()); err != nil {
		t.Fatalf("write image: %v", err)
	}

	if err := f.Close(); err != nil {
		t.Fatalf("close %s: %v", path, err)
	}
}

func prepareOutputDir(t *testing.T) {
	outputDir := filepath.Join(wd, "testoutput")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("create %s: %v", outputDir, err)
	}
}
