package fw

import (
	"os"
	"sync"
)

type FastWriter struct {
	fileName            string
	bufSize             int
	maxHistoryFileCount int
	fileMux             sync.Mutex
	f                   *os.File
	currentFileSize     int64
	mux                 sync.Mutex
	notifyWriterChan    chan struct{}
	bufAvailableCond    *sync.Cond
	buf                 []byte
	current             int
	writerWorking       bool
}

type FastWriterConfig struct {
	FileName            string
	BufferSize          int
	MaxHistoryFileCount int
}

func (conf *FastWriterConfig) init() error {
	if conf.FileName == "" {
		conf.FileName = "server.log"
	}
	if conf.BufferSize == 0 {
		conf.BufferSize = 512 * 1024
	}
	if conf.MaxHistoryFileCount == 0 {
		conf.MaxHistoryFileCount = 5
	}
	return nil
}

func NewFastWriter(conf *FastWriterConfig) (*FastWriter, error) {
	if err := conf.init(); err != nil {
		return nil, err
	}
	w := &FastWriter{
		fileName:            conf.FileName,
		bufSize:             conf.BufferSize,
		maxHistoryFileCount: conf.MaxHistoryFileCount,
		notifyWriterChan:    make(chan struct{}, 1),
		buf:                 make([]byte, conf.BufferSize),
	}
	w.bufAvailableCond = sync.NewCond(&w.mux)

	var err error
	w.f, err = os.OpenFile(w.fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	fStat, err := w.f.Stat()
	if err != nil {
		return nil, err
	}
	w.currentFileSize = fStat.Size()
	go w.writeLoop()
	return w, nil
}

func (w *FastWriter) writeLoop() {
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

func (w *FastWriter) Write(b []byte) (n int, err error) {
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

func (w *FastWriter) writeToFile(b []byte) (n int, err error) {
	w.fileMux.Lock()
	defer w.fileMux.Unlock()
	n, err = w.f.Write(b)
	w.currentFileSize += int64(n)
	return n, err
}

func (w *FastWriter) Close() error {
	w.mux.Lock()
	defer w.mux.Unlock()
	return w.f.Close()
}
