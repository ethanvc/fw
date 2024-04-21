package fw

import (
	"os"
	"sync"
)

type FastWriter struct {
	f         *os.File
	mux       sync.Mutex
	buf       []byte
	cacheSize int
	current   int
}

func NewFastWriter(fileName string) (*FastWriter, error) {
	w := &FastWriter{
		cacheSize: 1024,
	}
	err := w.init(fileName)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *FastWriter) init(fileName string) error {
	var err error
	w.f, err = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	w.buf = make([]byte, w.cacheSize)
	return nil
}

func (w *FastWriter) Write(p []byte) (n int, err error) {
	w.mux.Lock()
	defer w.mux.Unlock()
	if w.current+len(p) > w.cacheSize {
		n, err = w.f.Write(p)
		if err != nil {
			return n, err
		}
		w.current = 0
	}
	copy(w.buf[w.current:], p)
	w.current += len(p)
	return len(p), nil
}

func (w *FastWriter) Close() error {
	w.mux.Lock()
	defer w.mux.Unlock()
	if w.f == nil {
		return nil
	}
	if w.current > 0 {
		if _, err := w.f.Write(w.buf[:w.current]); err != nil {
			return err
		}
		w.current = 0
	}
	if err := w.f.Close(); err != nil {
		return err
	}
	w.f = nil
	return nil
}
