package fw

import (
	"bytes"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"testing"
)

func Benchmark_Lumberjack(b *testing.B) {
	w := &lumberjack.Logger{
		Filename: "lumberjack.test.log",
	}
	defer w.Close()
	// defer os.Remove(w.Filename)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			WriteLog(w)
		}
	})
}

func Benchmark_FileWriter(b *testing.B) {
	const fileName = "file_writer.test.log"
	w, err := NewFileWriter(fileName)
	if err != nil {
		panic(err)
	}
	defer w.Close()
	// defer os.Remove(fileName)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			WriteLog(w)
		}
	})
}

func Benchmark_MemoryMapWriter(b *testing.B) {
	const fileName = "memory_map_writer.test.log"
	w, err := NewMemoryMapWriter(fileName, 0)
	if err != nil {
		panic(err)
	}
	defer w.Close()
	// defer os.Remove(fileName)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			WriteLog(w)
		}
	})
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

func WriteLog(w io.Writer) {
	n, err := w.Write(testLogBuf)
	if err != nil {
		panic(err)
	}
	if n != len(testLogBuf) {
		panic("SizeNotMatch")
	}
}
