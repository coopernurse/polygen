package polygenlib

import (
	"os"
	"errors"
	"fmt"
	"bytes"
	"time"
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

type StrBuf struct {
    commentDelim string
	b *bytes.Buffer
}

func (sb StrBuf) prelude() {
	sb.c("DO NOT EDIT THIS FILE - Generated by polygen")
	sb.c("More info: https://github.com/coopernurse/polygen")
	sb.c(fmt.Sprintf("Created: %s", time.Now().Format(time.RFC1123)))
}

func (sb StrBuf) w(s string) {
	sb.b.WriteString(s)
	sb.b.WriteString("\n")
}

func (sb StrBuf) f(s string, args ...interface{}) {
	sb.w(fmt.Sprintf(s, args...))
}

func (sb StrBuf) c(s string) {
	sb.f("%s %s", sb.commentDelim, s)
}

func (sb StrBuf) blank() {
	sb.w("")
}

func (sb StrBuf) raw(s string) {
	sb.b.WriteString(s)
}

func (sb StrBuf) fraw(s string, args ...interface{}) {
	sb.raw(fmt.Sprintf(s, args...))
}

func NewStrBuf(commentDelim string) *StrBuf {
	return &StrBuf{commentDelim, &bytes.Buffer{}}
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

