package fw

import (
	"io"
	"strings"
)

var testLogBuf = []byte(strings.Repeat("a", 100))

func WriteLog(w io.Writer) {
	w.Write(testLogBuf)
}
