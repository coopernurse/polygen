package polygenlib

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestJsGenerator(t *testing.T) {
	pkg, err := Parse("test.go", example1)
	if err != nil {
		t.Fatal(err)
	}

	gen := JsGenerator{}
	files := gen.GenFiles(pkg)

	dir, err := ioutil.TempDir(os.TempDir(), "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	err = WriteFiles(dir, files)
	if err != nil {
		t.Fatal(err)
	}

}
