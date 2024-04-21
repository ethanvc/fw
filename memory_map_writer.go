package fw

import (
	"errors"
	"github.com/edsrzf/mmap-go"
	"os"
	"sync"
)

type MemoryMapWriter struct {
	f               *os.File
	blockSize       int
	mux             sync.Mutex
	currentFileSize int64
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
	err := w.init(fileName)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *MemoryMapWriter) init(fileName string) error {
	var err error
	w.f, err = os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	fStat, err := w.f.Stat()
	if err != nil {
		return err
	}
	w.currentFileSize = fStat.Size()
	return nil
}

func (w *MemoryMapWriter) Write(p []byte) (n int, err error) {
	w.mux.Lock()
	defer w.mux.Unlock()
	if err := w.reserveForNBytes(len(p)); err != nil {
		return 0, err
	}
	copy(w.block[w.current:], p)
	w.current += len(p)
	return len(p), nil
}

func (w *MemoryMapWriter) reserveForNBytes(n int) error {
	if len(w.block) >= n+w.current {
		return nil
	}
	w.block.Unmap()
	w.current = 0
	var err error
	err = w.f.Truncate(w.currentFileSize + int64(w.blockSize))
	if err != nil {
		return err
	}
	w.block, err = mmap.MapRegion(w.f, w.blockSize, mmap.RDWR, 0, w.currentFileSize)
	if err != nil {
		return err
	}
	w.currentFileSize += int64(w.blockSize)
	return nil
}

func (w *MemoryMapWriter) Close() error {
	w.mux.Lock()
	defer w.mux.Unlock()
	if w.f == nil {
		return nil
	}
	err := w.block.Unmap()
	if err != nil {
		return errors.Join(errors.New("UnmapErr"), err)
	}
	err = w.f.Truncate(w.currentFileSize - int64(w.blockSize-w.current))
	if err != nil {
		return errors.Join(errors.New("TruncateErr"), err)
	}
	err = w.f.Close()
	if err != nil {
		return errors.Join(errors.New("CloseErr"), err)
	}
	w.f = nil
	return nil
}
