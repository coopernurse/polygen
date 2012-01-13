package polygenlib

import (
	"testing"
	"os"
	"os/exec"
	"io/ioutil"
)

func TestJavaGenerator(t *testing.T) {
	pkg, err := Parse(example1)
	if err != nil {
		t.Fatal(err)
	}
	
	gen := JavaGenerator{}
	files := gen.GenFiles(pkg)
	names := []string{"Result.java", "Person.java",
		"SampleService.java", "SampleServiceRPCServer.java",
		"SampleServiceRPCClient.java"}
	if len(files) != len(names) {
		t.Errorf("Expected %d files, got: %d", len(names), len(files))
	}

	for i := 0; i < len(names); i++ {
		if names[i] != files[i].Name {
			t.Errorf("Filename %s != %s", names[i], files[i].Name)
		}
	}

	dir, err := ioutil.TempDir(os.TempDir(), "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	err = WriteFiles(dir, files)
	if err != nil {
		t.Fatal(err)
	}

	flist, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(flist) != len(names) {
		t.Errorf("%d != %d files in dir: %s", len(flist), len(names), dir)
	}

	args := []string{"-d", dir}
	for i := 0; i < len(files); i++ {
		args = append(args, files[i].FullPath(dir))
	}
	cmd := exec.Command("javac", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Error running command: javac %s\n%s\n%s", args, err, out)
	}
}