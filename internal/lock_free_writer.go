package internal

import (
	"sync/atomic"
	"unsafe"
)

type LockFreeWriter struct {
	block atomic.Pointer[bufferBlock]
}

func NewLockFreeWriter() *LockFreeWriter {
	w := &LockFreeWriter{}
	w.block.Store(newBufferBlock(1024 * 1024 * 5))
	return w
}

func (w *LockFreeWriter) Write(p []byte) (n int, err error) {
	block := w.block.Load()
	_ = block
	return len(p), nil
}

type bufferBlock struct {
	buf   []byte
	index atomic.Int64
}

func newBufferBlock(bufSize int) *bufferBlock {
	block := &bufferBlock{}
	block.buf = make([]byte, bufSize)
	return block
}

func (b *bufferBlock) tryWrite(p []byte) bool {
	bufLen := len(p) + 4
	idx := b.index.Add(int64(bufLen))
	if int(idx) > len(b.buf) {
		b.index.Add(-int64(bufLen))
		return false
	}
	startIdx := int(idx) - bufLen
	copy(b.buf[startIdx+4:], p)
	lenPointer := (*int32)(unsafe.Pointer(&b.buf[startIdx]))
	atomic.StoreInt32(lenPointer, int32(bufLen))
	return true
}
