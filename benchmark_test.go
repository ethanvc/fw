package fw

import (
	"bufio"
	"bytes"
	"github.com/ethanvc/fw/internal"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"sync"
	"testing"
)

func Benchmark_Lumberjack(b *testing.B) {
	w := &lumberjack.Logger{
		Filename: "lumberjack.test.log",
	}
	benchWriter(b, w)
}

func Benchmark_FileWriter(b *testing.B) {
	const fileName = "file_writer.test.log"
	os.Remove(fileName)
	w, err := internal.NewFileWriter(fileName)
	require.NoError(b, err)
	benchWriter(b, w)
}

func Benchmark_MemoryMapWriter(b *testing.B) {
	const fileName = "memory_map_writer.test.log"
	os.Remove(fileName)
	w, err := internal.NewMemoryMapWriter(fileName, 0)
	require.NoError(b, err)
	benchWriter(b, w)
}

func Benchmark_BatchWriter(b *testing.B) {
	const fileName = "batch.test.log"
	os.Remove(fileName)
	w, err := internal.NewBatchWriter(fileName)
	require.NoError(b, err)
	benchWriter(b, w)
}

func Benchmark_MultiBufferWriter(b *testing.B) {
	const fileName = "multi_buffer.test.log"
	os.Remove(fileName)
	w, err := internal.NewMultiBufferWriter(fileName)
	require.NoError(b, err)
	benchWriter(b, w)
}

func Benchmark_FastWriter(b *testing.B) {
	const fileName = "fast.test.log"
	os.Remove(fileName)
	w, err := NewFastWriter(&FastWriterConfig{
		Writer: &lumberjack.Logger{
			Filename: fileName,
		},
	})
	require.NoError(b, err)
	benchWriter(b, w)
}

func Benchmark_ZapBufferedWriter(b *testing.B) {
	const fileName = "zapbuffered.test.log"
	os.Remove(fileName)
	lumL := &lumberjack.Logger{
		Filename: fileName,
	}
	ws := zapcore.AddSync(lumL)
	w := &zapcore.BufferedWriteSyncer{
		WS: ws,
	}
	benchWriter(b, w)
}

func Benchmark_NopWriter(b *testing.B) {
	w := internal.NewNopWriter()
	benchWriter(b, w)
}

func Benchmark_BufioWriter(b *testing.B) {
	w := bufio.NewWriter(&lumberjack.Logger{})
	benchWriter(b, newSequenceWriteCloser(w))
}

type nopWriteCloser struct {
	mux sync.Mutex
	io.Writer
}

func newSequenceWriteCloser(w io.Writer) io.WriteCloser {
	return &nopWriteCloser{
		Writer: w,
	}
}

func (w *nopWriteCloser) Write(b []byte) (n int, err error) {
	w.mux.Lock()
	defer w.mux.Unlock()
	return w.Writer.Write(b)
}

func (w *nopWriteCloser) Close() error {
	return nil
}

func benchWriter(b *testing.B, w io.Writer) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n, err := w.Write(testLogBuf)
			require.NoError(b, err)
			require.Equal(b, len(testLogBuf), n)
		}
	})
	if closer, ok := w.(io.Closer); ok {
		require.NoError(b, closer.Close())
	}
}

var testLogBuf = generateTestData()

func generateTestData() []byte {
	var buf bytes.Buffer
	const maxCount = 26
	for i := 0; i < 26; i++ {
		if buf.Len() >= maxCount {
			break
		}
		buf.WriteRune(rune('A' + (i % 26)))
	}
	buf.WriteRune('\n')
	return buf.Bytes()
}
