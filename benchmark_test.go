package fw

import (
	"bytes"
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
	w, err := NewFileWriter(fileName)
	require.NoError(b, err)
	benchWriter(b, w)
}

func Benchmark_MemoryMapWriter(b *testing.B) {
	const fileName = "memory_map_writer.test.log"
	os.Remove(fileName)
	w, err := NewMemoryMapWriter(fileName, 0)
	require.NoError(b, err)
	benchWriter(b, w)
}

func Benchmark_FastWriter(b *testing.B) {
	const fileName = "fast.test.log"
	os.Remove(fileName)
	w, err := NewFastWriter(fileName)
	require.NoError(b, err)
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
