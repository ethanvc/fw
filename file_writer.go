package fw

import (
	"os"
)

type FileWriter struct {
	f *os.File
}

func NewFileWriter(fileName string) (*FileWriter, error) {
	w := &FileWriter{}
	err := w.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	return w, nil
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
