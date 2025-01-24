package daikin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeDecode(t *testing.T) {
	secretKey := []byte("YELLOW SUBMARINE")
	encodedData, err := encodeData(secretKey, []byte(`{"port1":{"power":1}}`))
	require.NoError(t, err)
	decodedData, err := decodeData(secretKey, encodedData, "")
	require.NoError(t, err)
	require.Equal(t, []byte(`{"port1":{"power":1}}BZ`), decodedData)
}
