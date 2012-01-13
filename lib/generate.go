package polygenlib

import (
	"os"
	"errors"
	"fmt"
)

type CodeGenerator interface {
	GenFiles(p *Package) []File
}

type File struct {
	Name string
	Contents []byte
}

func (f File) FullPath(dir string) string {
	return fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), f.Name)
}

func (f File) WriteTo(dir string) error {
	out, err := os.Create(f.FullPath(dir))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.Write(f.Contents)
	return err
}

func WriteFiles(dir string, files []File) error {
	fi, err := os.Stat(dir)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return errors.New(fmt.Sprintf("Invalid directory: %s", dir))
	}

	for i := 0; i < len(files); i++ {
		err := files[i].WriteTo(dir)
		if err != nil {
			return err
		}
	}

	return nil
}