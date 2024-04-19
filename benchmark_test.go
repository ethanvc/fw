package fw

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strings"
	"testing"
)

func Benchmark_FileWriter(b *testing.B) {
	w := NewFileWriter()
	fileName := "file_writer.test.log"
	err := w.OpenFile(fileName)
	if err != nil {
		panic(err)
	}
	defer w.Close()
	defer os.Remove(fileName)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			WriteLog(w)
		}
	})
}

func Benchmark_Lumberjack(b *testing.B) {
	w := &lumberjack.Logger{
		Filename: "lumberjack.test.log",
	}
	defer w.Close()
	defer os.Remove(w.Filename)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			WriteLog(w)
		}
	})
}

var testLogBuf = []byte(strings.Repeat("a", 100))

func WriteLog(w io.Writer) {
	w.Write(testLogBuf)
}
