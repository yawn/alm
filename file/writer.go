package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Writer struct {
	directory string
	known     map[string]interface{}
}

func New(directory string) *Writer {

	return &Writer{
		directory: directory,
		known:     make(map[string]interface{}),
	}

}

func (w *Writer) Add(parts ...string) (io.WriteCloser, error) {

	file := fmt.Sprintf("%s.log", filepath.Join(w.directory, strings.Join(parts, "-")))

	if _, ok := w.known[file]; !ok {

		w.known[file] = struct{}{}
		os.Remove(file)

	}

	return os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

}
