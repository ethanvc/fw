package fw

import (
	"os"
)

type FileWriter struct {
	f *os.File
}

func NewFileWriter() *FileWriter {
	w := &FileWriter{}
	return w
}

func (w *FileWriter) OpenFile(fileName string) error {
	if w.f != nil {
		w.f.Close()
	}
	var err error
	w.f, err = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	return w.f.Write(p)
}

func (w *FileWriter) Close() error {
	return w.f.Close()
}
