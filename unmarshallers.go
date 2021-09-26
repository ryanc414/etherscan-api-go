package etherscan

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	if bytes.Equal(data, []byte("\"0x\"")) {
		*u = 0
		return nil
	}

	var val hexutil.Uint64
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	*u = hexUint(val)
	return nil
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

func unmarshalResponse(data []byte, v interface{}) error {
	rspType := reflect.TypeOf(v)
	if rspType.Kind() != reflect.Ptr {
		return errors.New("value must be a pointer")
	}

	rspStructType := rspType.Elem()
	switch rspStructType.Kind() {
	case reflect.Struct:
		return unmarshalStructRsp(data, v)

	case reflect.Slice:
		if rspStructType.Elem().Kind() != reflect.Struct {
			return errors.New("only slices of structs are allowed")
		}

		return unmarshalSliceRsp(data, v)

	default:
		return errors.New("value must be a pointer to struct or slice")
	}
}

func unmarshalSliceRsp(data []byte, v interface{}) error {
	val := reflect.ValueOf(v).Elem()

	var rawSlice []json.RawMessage
	if err := json.Unmarshal(data, &rawSlice); err != nil {
		return errors.Wrap(err, "while unmarshalling as slice")
	}

	slice := reflect.MakeSlice(val.Type(), len(rawSlice), len(rawSlice))

	for i := range rawSlice {
		v := slice.Index(i).Addr().Interface()

		if err := unmarshalStructRsp(rawSlice[i], v); err != nil {
			return err
		}
	}

	val.Set(slice)

	return nil
}

func unmarshalStructRsp(data []byte, v interface{}) error {
	var rspMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rspMap); err != nil {
		return errors.Wrap(err, "while unmarshalling as map")
	}

	val := reflect.ValueOf(v).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		name := getFieldName(val.Type().Field(i))
		fieldData := rspMap[name]
		if err := setFieldValue(field, fieldData); err != nil {
			return err
		}
	}

	return nil
}

func getFieldName(field reflect.StructField) string {
	return strings.ToLower(field.Name)
}

func setFieldValue(field reflect.Value, data []byte) error {
	iField := field.Interface()
	if _, ok := iField.(*big.Int); ok {
		var res bigInt
		if err := json.Unmarshal(data, &res); err != nil {
			return errors.Wrap(err, "while unmarshalling as bigInt")
		}

		field.Set(reflect.ValueOf(res.unwrap()))
		return nil
	}

	if field.Kind() == reflect.Ptr {
		if err := json.Unmarshal(data, iField); err != nil {
			return errors.Wrap(err, "while unmarshalling as pointer")
		}

		return nil
	}

	if err := json.Unmarshal(data, field.Addr().Interface()); err != nil {
		return errors.Wrap(err, "while unmarshalling as value")
	}

	return nil
}
