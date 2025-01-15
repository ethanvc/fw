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
	block.tryWrite(p)
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
	bufIndex := b.tryGetBuffer(len(p))
	if bufIndex == -1 {
		return false
	}
	copy(b.buf[bufIndex+4:], p)
	atomic.StoreInt32((*int32)(unsafe.Pointer(&b.buf[bufIndex])), int32(4+len(p)))
	return true
}

func (b *bufferBlock) tryGetBuffer(contentLen int) int {
	// 4 for content len, 3 for memory alignment
	realBufLen := contentLen + 4 + 3
	realIndex := int(b.index.Add(int64(realBufLen)))
	if realIndex > len(b.buf) {
		b.index.Add(-int64(realBufLen))
		return -1
	}
	realIndex -= realBufLen
	alignedIndex := realIndex + 3 - (realIndex+3)%4
	return alignedIndex
}
