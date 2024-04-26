package fw

import (
	"bytes"
	"github.com/ethanvc/fw/internal"
	"github.com/stretchr/testify/require"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
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
		Writer: &lumberjack.Logger{},
	})
	require.NoError(b, err)
	benchWriter(b, w)
}

func Benchmark_NopWriter(b *testing.B) {
	w := internal.NewNopWriter()
	benchWriter(b, w)
}

func benchWriter(b *testing.B, w io.WriteCloser) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n, err := w.Write(testLogBuf)
			require.NoError(b, err)
			require.Equal(b, len(testLogBuf), n)
		}
	})
	require.NoError(b, w.Close())
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
