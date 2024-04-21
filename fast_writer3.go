package fw

import (
	"os"
	"runtime"
	"sync"
	"time"
)

type FastWriter3 struct {
	f       *os.File
	mux     sync.Mutex
	bufSize int
	buf     []byte
	current int
}

func NewFastWriter3(fileName string) (*FastWriter3, error) {
	w := &FastWriter3{
		bufSize: 1024 * 1024,
	}
	err := w.init(fileName)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *FastWriter3) init(fileName string) error {
	w.buf = make([]byte, w.bufSize)
	var err error
	w.f, err = os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	go w.writeLoop()
	return nil
}

func (w *FastWriter3) writeLoop() {
	buf := make([]byte, w.bufSize)
	for {
		w.mux.Lock()
		if w.current == 0 {
			w.mux.Unlock()
			time.Sleep(time.Millisecond)
			continue
		}
		buf, w.buf = w.buf, buf
		contentSize := w.current
		w.current = 0
		w.mux.Unlock()
		_, err := w.f.Write(buf[0:contentSize])
		if err != nil {
			break
		}
	}
}

func (w *FastWriter3) Write(p []byte) (n int, err error) {
	for {
		w.mux.Lock()
		if len(p)+w.current > len(w.buf) {
			w.mux.Unlock()
			runtime.Gosched()
			continue
		}
		copy(w.buf[w.current:], p)
		w.current += len(p)
		w.mux.Unlock()
		break
	}
	return len(p), nil
}

func (w *FastWriter3) Close() error {
	w.mux.Lock()
	defer w.mux.Unlock()
	return w.f.Close()
}
