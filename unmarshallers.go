package etherscan

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	return (*big.Int)(b)
}

type hexBigInt big.Int

func (b *hexBigInt) UnmarshalJSON(data []byte) error {
	var hex string
	if err := json.Unmarshal(data, &hex); err != nil {
		return err
	}

	val := new(big.Int).SetBytes(common.Hex2Bytes(hex))
	*b = hexBigInt(*val)
	return nil
}

func (b *hexBigInt) unwrap() *big.Int {
	return (*big.Int)(b)
}

type uintStr uint64

func (u *uintStr) UnmarshalJSON(data []byte) error {
	var rawStr string
	if err := json.Unmarshal(data, &rawStr); err != nil {
		return err
	}
	if rawStr == "" {
		return nil
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

type hexUint uint64

func (u *hexUint) UnmarshalJSON(data []byte) error {
	var hex string
	if err := json.Unmarshal(data, &hex); err != nil {
		return err
	}

	buf := bytes.NewBuffer(common.Hex2BytesFixed(hex, 64))

	var val uint64
	if err := binary.Read(buf, binary.BigEndian, &val); err != nil {
		return err
	}

	*u = hexUint(val)
	return nil
}

func (u hexUint) unwrap() uint64 {
	return uint64(u)
}

type unixTimestamp time.Time

func (t *unixTimestamp) UnmarshalJSON(data []byte) error {
	var rawStr string
	if err := json.Unmarshal(data, &rawStr); err != nil {
		return err
	}

	unixSeconds, err := strconv.ParseInt(rawStr, 10, 64)
	if err != nil {
		return err
	}

	*t = unixTimestamp(time.Unix(unixSeconds, 0))
	return nil
}

func (t unixTimestamp) unwrap() time.Time {
	return time.Time(t)
}

type hexTimestamp time.Time

func (t *hexTimestamp) UnmarshalJSON(data []byte) error {
	var hex string
	if err := json.Unmarshal(data, &hex); err != nil {
		return err
	}

	buf := bytes.NewBuffer(common.Hex2BytesFixed(hex, 64))

	var unixSeconds uint64
	if err := binary.Read(buf, binary.BigEndian, &unixSeconds); err != nil {
		return err
	}

	*t = hexTimestamp(time.Unix(int64(unixSeconds), 0))
	return nil

}

func (t hexTimestamp) unwrap() time.Time {
	return time.Time(t)
}

type floatStr float64

func (f *floatStr) UnmarshalJSON(data []byte) error {
	var rawStr string
	if err := json.Unmarshal(data, &rawStr); err != nil {
		return err
	}
	if rawStr == "" {
		return nil
	}

	val, err := strconv.ParseFloat(rawStr, 64)
	if err != nil {
		return err
	}

	*f = floatStr(val)
	return nil
}

func (f floatStr) unwrap() float64 {
	return float64(f)
}
