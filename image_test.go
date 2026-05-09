package escposimg

import (
	"image/color"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadImageDetails_PBM(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(thisFile)
	path := filepath.Join(dir, "testdata", "tiny.pbm")

	res, err := LoadImageDetails(path)
	if err != nil {
		t.Fatal(err)
	}
	if res.Format != "pbm" {
		t.Fatalf("format: got %q want pbm", res.Format)
	}
	if !res.SkipDithering {
		t.Fatal("SkipDithering: want true for PBM")
	}

	b := res.Image.Bounds()
	if b.Dx() != 2 || b.Dy() != 2 {
		t.Fatalf("bounds: got %dx%d want 2x2", b.Dx(), b.Dy())
	}

	gray := color.GrayModel.Convert(res.Image.At(b.Min.X, b.Min.Y)).(color.Gray)
	if gray.Y >= 128 {
		t.Fatalf("top-left pixel: got Y=%d want black (<128)", gray.Y)
	}
}

func TestLoadImage_Wrapper(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(thisFile)
	path := filepath.Join(dir, "testdata", "tiny.pbm")

	img, err := LoadImage(path)
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Empty() {
		t.Fatal("empty image")
	}
}
