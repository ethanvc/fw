package fw

import (
	"errors"
	"io"
	"os"
	"sync"
)

type FastWriter struct {
	bufSize          int
	writerMux        sync.Mutex
	writer           io.Writer
	mux              sync.Mutex
	notifyWriterChan chan struct{}
	bufAvailableCond *sync.Cond
	buf              []byte
	current          int
	writerWorking    bool
	closed           bool
}

type FastWriterConfig struct {
	Writer     io.Writer
	BufferSize int
}

func (conf *FastWriterConfig) init() error {
	if conf.Writer == nil {
		conf.Writer = os.Stdout
	}
	if conf.BufferSize == 0 {
		conf.BufferSize = 512 * 1024
	}
	return nil
}

func NewFastWriter(conf *FastWriterConfig) (*FastWriter, error) {
	if err := conf.init(); err != nil {
		return nil, err
	}
	w := &FastWriter{
		writer:           conf.Writer,
		bufSize:          conf.BufferSize,
		notifyWriterChan: make(chan struct{}, 1),
		buf:              make([]byte, conf.BufferSize),
	}
	w.bufAvailableCond = sync.NewCond(&w.mux)
	go w.writeLoop()
	return w, nil
}

func (w *FastWriter) writeLoop() {
	buf := make([]byte, w.bufSize)
	for {
		var wrote bool
		buf, wrote = w.writeOnce(buf)
		if wrote {
			continue
		}
		select {
		case <-w.notifyWriterChan:
		}
		if w.closed {
			break
		}
	}
}

func (w *FastWriter) writeOnce(buf []byte) (freeBuf []byte, wrote bool) {
	w.writerMux.Lock()
	defer w.writerMux.Unlock()
	buf, n := w.exchangeBuffer(buf)
	if n == 0 {
		return buf, false
	}
	w.bufAvailableCond.Broadcast()
	w.writer.Write(buf[0:n])
	return buf, true
}

func (w *FastWriter) exchangeBuffer(buf []byte) ([]byte, int) {
	w.mux.Lock()
	defer w.mux.Unlock()
	if w.current == 0 {
		w.writerWorking = false
		return buf, 0
	}
	w.writerWorking = true
	n := w.current
	w.current = 0
	buf, w.buf = w.buf, buf
	return buf, n
}

func (w *FastWriter) Write(b []byte) (n int, err error) {
	l := len(b)
	if l > w.bufSize {
		return w.writeDirect(b)
	}
	w.mux.Lock()
	for l+w.current > len(w.buf) && !w.closed {
		w.bufAvailableCond.Wait()
	}
	if w.closed {
		w.mux.Unlock()
		return 0, errors.New("AlreadyClosed")
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

func (w *FastWriter) writeDirect(b []byte) (n int, err error) {
	w.writerMux.Lock()
	defer w.writerMux.Unlock()
	return w.writer.Write(b)
}

func (w *FastWriter) Flush() error {
	buf := make([]byte, w.bufSize)
	w.writeOnce(buf)
	w.flushWriter()
	return nil
}

func (w *FastWriter) flushWriter() {
	if f, _ := w.writer.(interface{ Flush() error }); f != nil {
		f.Flush()
		return
	}
	if f, _ := w.writer.(interface{ Sync() error }); f != nil {
		f.Sync()
		return
	}
}

func (w *FastWriter) Close() error {
	w.mux.Lock()
	w.closed = true
	w.mux.Unlock()
	select {
	case w.notifyWriterChan <- struct{}{}:
	default:
	}

	return w.Flush()
}
