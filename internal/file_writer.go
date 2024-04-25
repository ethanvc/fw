package internal

import (
	"os"
	"sync"
)

type FileWriter struct {
	f   *os.File
	mux sync.Mutex
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
	w.Close()
	var err error
	w.f, err = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	w.mux.Lock()
	defer w.mux.Unlock()
	return w.f.Write(p)
}

func (w *FileWriter) Close() error {
	var f *os.File
	w.mux.Lock()
	f = w.f
	w.f = nil
	w.mux.Unlock()
	if f == nil {
		return nil
	}
	return f.Close()
}
