package daikin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_calculateCRC16(t *testing.T) {
	t.Parallel()
	r := calculateCRC16([]byte("test"))
	require.Equal(t, -8135, r)
}
