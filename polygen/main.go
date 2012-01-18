package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"path/filepath"
	polygen "github.com/coopernurse/polygen/lib"
)

func generate(p *polygen.Package, g polygen.CodeGenerator, dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	log.Printf("Writing to directory: %s", dir)

	files := g.GenFiles(p)
	for i := 0; i < len(files); i++ {
		log.Printf("   Writing: %s", files[i].Name)
		err := files[i].WriteTo(dir)
		if err != nil {
			return err
		}
	}

	return nil
}

func usage() string {
	b := bytes.Buffer{}
	b.WriteString("usage: polygen [options] idlfile\n")
	flag.VisitAll(func(f *flag.Flag) {
		b.WriteString(fmt.Sprintf("  -%s=%s: %s)\n", f.Name, f.DefValue, f.Usage))
	})
	return b.String()
}

func main() {
	var dir string
	var clean bool
	flag.StringVar(&dir, "dir", ".", "directory to write files to")
	flag.BoolVar(&clean, "c", false, "delete existing files from dirs")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal(usage())
	}

	fname := args[0]
	log.Printf("Reading IDL from: %s", fname)

	idl, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	pkg, err := polygen.Parse(string(idl))
	if err != nil {
		log.Fatal(err)
	}	

	generators := make(map[string]polygen.CodeGenerator)
	generators["java"] = polygen.JavaGenerator{}
	generators["js"] = polygen.JsGenerator{}

	for subdir, gen := range(generators) {
		dest := filepath.Join(dir, subdir)
		if clean {
			log.Printf("Cleaning: %s", dest)
			err := os.RemoveAll(dest)
			if err != nil {
				log.Fatal(err)
			}
		}
		err = generate(pkg, gen, dest)
		if err != nil {
			log.Fatal(err)
		}
	}
}