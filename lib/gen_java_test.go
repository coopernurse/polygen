package polygenlib

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestJavaGenerator(t *testing.T) {
	pkg, err := Parse("java_test.go", example1)
	if err != nil {
		t.Fatal(err)
	}

	gen := JavaGenerator{}
	files := gen.GenFiles(pkg)
	names := []string{"RPCException.java", "RPCError.java",
		"Result.java", "Person.java",
		"SampleService.java", "SampleServiceDispatcher.java",
		"SampleServiceHttpServer.java",
		"SampleServiceClient.java", "SampleServiceTypes.java"}
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

	//dir := "/Users/james/go/src/github.com/coopernurse/polygen/lib/test"

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

	args := []string{"-cp", "deps/jackson-mapper-lgpl-1.9.3.jar:deps/jackson-core-lgpl-1.9.3.jar", "-d", dir}
	for i := 0; i < len(files); i++ {
		args = append(args, files[i].FullPath(dir))
	}
	cmd := exec.Command("javac", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Error running command: javac %s\n%s\n%s", args, err, out)
	}
}
