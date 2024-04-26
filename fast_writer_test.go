package fw

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"runtime"
	"sync"
	"testing"
)

func TestFastWriter_ConcurrentWrite(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w, err := NewFastWriter(&FastWriterConfig{
		Writer: buf,
	})
	require.NoError(t, err)
	var wg sync.WaitGroup
	concurrentCount := runtime.NumCPU() * 3
	wg.Add(concurrentCount)
	const contentLen = 10000
	for i := 0; i < concurrentCount; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 10000; i++ {
				w.Write([]byte("a"))
			}
		}()
	}
	wg.Wait()
	w.Flush()
	require.Equal(t, contentLen*concurrentCount, buf.Len())
}
