package fw

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMemoryMapWriter_Basic(t *testing.T) {
	const fileName = "test.log"
	w, err := NewMemoryMapWriter(fileName, 0)
	require.NoError(t, err)
	_, err = w.Write([]byte("test"))
	require.NoError(t, err)
	err = w.Close()
	require.NoError(t, err)
	os.Remove(fileName)
}
