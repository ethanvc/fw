package fw

import (
	"github.com/edsrzf/mmap-go"
	"io"
	"os"
	"sync"
)

type MemoryMapWriter struct {
	f               *os.File
	blockSize       int
	mux             sync.Mutex
	CurrentFileSize int64
	block           mmap.MMap
	current         int
}

func NewMemoryMapWriter(fileName string, blockSize int) (*MemoryMapWriter, error) {
	if blockSize == 0 {
		blockSize = 1024 * 64
	}
	w := &MemoryMapWriter{
		blockSize: blockSize,
	}
	err := w.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *MemoryMapWriter) OpenFile(fileName string) error {
	w.Close()
	var err error
	w.f, err = os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	fStat, err := w.f.Stat()
	if err != nil {
		return err
	}
	w.CurrentFileSize = fStat.Size()
	return nil
}

func (w *MemoryMapWriter) Write(p []byte) (n int, err error) {
	w.mux.Lock()
	defer w.mux.Unlock()
	if err := w.reserveForNBytes(len(p)); err != nil {
		return 0, err
	}
	copy(w.block[w.current:], p)
	return len(p), nil
}

func (w *MemoryMapWriter) reserveForNBytes(n int) error {
	if len(w.block) >= n+w.current {
		return nil
	}
	w.block.Unmap()
	w.current = 0
	var err error
	_, err = w.f.Seek(w.CurrentFileSize+int64(w.blockSize), io.SeekStart)
	if err != nil {
		return err
	}
	w.block, err = mmap.MapRegion(w.f, w.blockSize, mmap.RDWR, 0, w.CurrentFileSize)
	if err != nil {
		return err
	}
	return nil
}

func (w *MemoryMapWriter) Close() error {
	w.mux.Lock()
	defer w.mux.Unlock()
	err := w.block.Unmap()
	if err != nil {
		return err
	}
	w.block = nil

	if w.f != nil {
		err = w.f.Close()
		if err != nil {
			return err
		}
		w.f = nil
	}
	return nil
}
