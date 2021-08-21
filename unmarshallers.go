package etherscan

import (
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/pkg/errors"
)

type bigInt big.Int

func (b *bigInt) UnmarshalJSON(data []byte) error {
	var intStr string
	if err := json.Unmarshal(data, &intStr); err != nil {
		return err
	}

	val, ok := new(big.Int).SetString(intStr, 10)
	if !ok {
		return errors.Errorf("cannot parse %s as big.Int", intStr)
	}

	*b = bigInt(*val)
	return nil
}

func (b *bigInt) unwrap() *big.Int {
	v := big.Int(*b)
	return &v
}

type uintStr uint64

func (u *uintStr) UnmarshalJSON(data []byte) error {
	var rawStr string
	if err := json.Unmarshal(data, &rawStr); err != nil {
		return err
	}

	val, err := strconv.ParseUint(rawStr, 10, 64)
	if err != nil {
		return err
	}

	*u = uintStr(val)
	return nil
}

func (u uintStr) unwrap() uint64 {
	return uint64(u)
}
