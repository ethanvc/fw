package fw

import (
	"os"
	"sync"
)

type FastWriter3 struct {
	fileMux          sync.Mutex
	f                *os.File
	mux              sync.Mutex
	notifyWriterChan chan struct{}
	bufAvailableCond *sync.Cond
	bufSize          int
	buf              []byte
	current          int
	writerWorking    bool
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
	w.bufAvailableCond = sync.NewCond(&w.mux)
	w.notifyWriterChan = make(chan struct{}, 1)
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
		for w.current == 0 {
			w.writerWorking = false
			w.mux.Unlock()
			select {
			case <-w.notifyWriterChan:
			}
			w.mux.Lock()
		}
		buf, w.buf = w.buf, buf
		contentSize := w.current
		w.current = 0
		w.writerWorking = true
		w.mux.Unlock()
		w.bufAvailableCond.Broadcast()
		w.writeToFile(buf[0:contentSize])
	}
}

func (w *FastWriter3) Write(b []byte) (n int, err error) {
	l := len(b)
	if l > w.bufSize {
		return w.writeToFile(b)
	}
	w.mux.Lock()
	for l+w.current > len(w.buf) {
		w.bufAvailableCond.Wait()
	}
	copy(w.buf[w.current:], b)
	w.current += l
	writerWorking := w.writerWorking
	w.mux.Unlock()
	if !writerWorking {
		select {
		case w.notifyWriterChan <- struct{}{}:
		default:
		}
	}
	return l, nil
}

func (w *FastWriter3) writeToFile(b []byte) (n int, err error) {
	w.fileMux.Lock()
	defer w.fileMux.Unlock()
	return w.f.Write(b)
}

func (w *FastWriter3) Close() error {
	w.mux.Lock()
	defer w.mux.Unlock()
	return w.f.Close()
}
