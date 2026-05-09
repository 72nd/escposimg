package escposimg

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestProcessImage_PBM(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(thisFile)
	imgPath := filepath.Join(dir, "testdata", "tiny.pbm")
	outPath := filepath.Join(t.TempDir(), "out.escpos")

	out, err := NewFileOutput(outPath)
	if err != nil {
		t.Fatal(err)
	}

	cfg := DefaultConfig()
	if err := ProcessImage(imgPath, cfg, out); err != nil {
		t.Fatal(err)
	}
	st, err := os.Stat(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if st.Size() == 0 {
		t.Fatal("expected non-empty ESC/POS output")
	}
}
