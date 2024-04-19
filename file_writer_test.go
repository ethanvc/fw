package fw

import (
	"os"
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
