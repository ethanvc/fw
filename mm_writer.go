package fw

import (
	"github.com/edsrzf/mmap-go"
	"os"
)

type MemoryMapWriter struct {
	f     *os.File
	block mmap.MMap
}
