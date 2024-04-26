package fw

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_getHistoryFileName(t *testing.T) {
	tim, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")
	require.NoError(t, err)
	require.Equal(t, "server-2006-01-02T15-04-05.000000000Z08-00.log", getHistoryFileName("server.log", tim))
}
