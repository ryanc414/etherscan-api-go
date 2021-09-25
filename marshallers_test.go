package etherscan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestMarshaller(t *testing.T) {
	req := BlockNumberAndIndex{
		Number: 123456,
		Index:  420,
	}

	res := marshalRequest(&req)

	expected := map[string]string{
		"tag":   "0x1e240",
		"index": "0x1a4",
	}
	assert.Equal(t, expected, res)
}
