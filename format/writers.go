package format

import (
	serrors "errors"
	"io"
	"io/ioutil"

	"github.com/atotto/clipboard"
)

type clipWriter struct{}

func (c clipWriter) Write(p []byte) (n int, err error) {
	if len(p) <= 0 {
		return 0, nil
	}
	data := string(p)
	if err := clipboard.WriteAll(data); err == nil {
		n = len(p)
	} else {
		err = serrors.New("Error writing to clipboard; " + err.Error())
	}
	return n, err
}

type fileWriter struct {
	filePath string
}

func NewFileWriter(filePath string) io.Writer {
	return &fileWriter{filePath}
}

func (w fileWriter) Write(p []byte) (n int, err error) {
	err = ioutil.WriteFile(w.filePath, p, 0664)
	if err == nil {
		n = len(p)
	}
	return n, err
}
