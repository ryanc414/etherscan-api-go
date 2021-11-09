package marshallers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type blockNumberAndIndex struct {
	Number uint64 `etherscan:"tag,hex"`
	Index  uint32 `etherscan:"index,hex"`
}

func TestRequestMarshaller(t *testing.T) {
	req := blockNumberAndIndex{
		Number: 123456,
		Index:  420,
	}

	res := MarshalRequest(&req)

	expected := map[string]string{
		"tag":   "0x1e240",
		"index": "0x1a4",
	}
	assert.Equal(t, expected, res)
}
