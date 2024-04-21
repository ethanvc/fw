package fw

import (
	"os"
)

type FastWriter2 struct {
	f       *os.File
	bufChan chan []byte
}

func NewFastWriter2(fileName string) (*FastWriter2, error) {
	w := &FastWriter2{
		bufChan: make(chan []byte, 1000),
	}
	err := w.init(fileName)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *FastWriter2) init(fileName string) error {
	var err error
	w.f, err = os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	go w.writeLoop()
	return nil
}

func (w *FastWriter2) Write(buf []byte) (n int, err error) {
	var newBuf []byte
	newBuf = append(newBuf, buf...)
	w.bufChan <- newBuf
	return len(buf), nil
}

func (w *FastWriter2) Close() error {
	return w.f.Close()
}

func (w *FastWriter2) writeLoop() {
	for buf := range w.bufChan {
		w.f.Write(buf)
	}
}
