package etherscan

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestMarshaller(t *testing.T) {
	req := BlockNumberAndIndex{
		Number: 123456,
		Index:  420,
	}

	res, err := marshalRequest(req)
	require.NoError(t, err)

	expected := map[string]string{
		"tag":   "123456",
		"index": "420",
	}
	assert.Equal(t, expected, res)
}
