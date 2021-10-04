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
	"github.com/rs/zerolog/log"
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

	rspVal := reflect.ValueOf(v).Elem()

	switch rspVal.Kind() {
	case reflect.Struct:
		return unmarshalStructRsp(data, rspVal)

	case reflect.Slice:
		if rspVal.Type().Elem().Kind() != reflect.Struct {
			return errors.New("only slices of structs are allowed")
		}

		return unmarshalSliceRsp(data, rspVal)

	default:
		return errors.New("value must be a pointer to struct or slice")
	}
}

func unmarshalSliceRsp(data []byte, v reflect.Value) error {
	var u json.Unmarshaler
	if reflect.PtrTo(v.Type().Elem()).Implements(reflect.TypeOf(&u).Elem()) {
		return json.Unmarshal(data, v.Addr().Interface())
	}

	var rawSlice []json.RawMessage
	if err := json.Unmarshal(data, &rawSlice); err != nil {
		return errors.Wrap(err, "while unmarshalling as slice")
	}

	slice := reflect.MakeSlice(v.Type(), len(rawSlice), len(rawSlice))

	for i := range rawSlice {
		el := slice.Index(i)

		if err := unmarshalStructRsp(rawSlice[i], el); err != nil {
			return err
		}
	}

	v.Set(slice)

	return nil
}

func unmarshalStructRsp(data []byte, v reflect.Value) error {
	var rspMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rspMap); err != nil {
		return errors.Wrap(err, "while unmarshalling as map")
	}

	fieldTypes := reflect.VisibleFields(v.Type())

	for i := range fieldTypes {
		fieldType := fieldTypes[i]
		if fieldType.Anonymous {
			continue
		}

		field := v.FieldByIndex(fieldType.Index)

		info := parseTag(fieldType)
		name := getFieldName(fieldType, &info)

		fieldData := rspMap[name]
		if len(fieldData) == 0 {
			log.Debug().Msgf("no field with name %s in response data", name)
			continue
		}

		if err := setFieldValue(field, fieldData, &info); err != nil {
			return errors.Wrapf(err, "while unmarshalling field %s", name)
		}
	}

	return nil
}

func getFieldName(field reflect.StructField, info *tagInfo) string {
	if info.name != "" {
		return info.name
	}

	return strings.ToLower(field.Name)
}

func setFieldValue(field reflect.Value, data []byte, info *tagInfo) error {
	if string(data) == "\"\"" {
		return nil
	}

	iField := field.Interface()
	if _, ok := iField.(*big.Int); ok {
		if info.hex {
			var res hexutil.Big
			if err := json.Unmarshal(data, &res); err != nil {
				return errors.Wrap(err, "while unmarshalling as hexutil.Big")
			}

			field.Set(reflect.ValueOf(res.ToInt()))
			return nil
		}

		var res bigInt
		if err := json.Unmarshal(data, &res); err != nil {
			return errors.Wrap(err, "while unmarshalling as bigInt")
		}

		field.Set(reflect.ValueOf(res.unwrap()))
		return nil
	}

	if _, ok := iField.(time.Time); ok {
		if info.hex {
			var res hexTimestamp
			if err := json.Unmarshal(data, &res); err != nil {
				return errors.Wrap(err, "while unmarshalling as unix timestamp")
			}

			field.Set(reflect.ValueOf(res.unwrap()))
			return nil
		}

		var res unixTimestamp
		if err := json.Unmarshal(data, &res); err != nil {
			return errors.Wrap(err, "while unmarshalling as unix timestamp")
		}

		field.Set(reflect.ValueOf(res.unwrap()))
		return nil
	}

	if _, ok := iField.([]byte); ok {
		if string(data) == "\"deprecated\"" {
			return nil
		}

		var res string
		if err := json.Unmarshal(data, &res); err != nil {
			return errors.Wrap(err, "while unmarshalling as string")
		}

		decoded, err := hexutil.Decode(res)
		if err != nil {
			return errors.Wrap(err, "while decoding as hex")
		}

		field.SetBytes(decoded)
		return nil
	}

	switch field.Kind() {
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		if info.num {
			if err := json.Unmarshal(data, field.Addr().Interface()); err != nil {
				return errors.Wrap(err, "while unmarshalling as uint")
			}

			return nil
		}

		if info.hex {
			var res hexUint
			if err := json.Unmarshal(data, &res); err != nil {
				return errors.Wrap(err, "while unmarshalling as hexUint")
			}

			field.SetUint(uint64(res))
			return nil
		}

		var res uintStr
		if err := json.Unmarshal(data, &res); err != nil {
			return errors.Wrap(err, "while unmarshalling as uintStr")
		}

		field.SetUint(res.unwrap())
		return nil

	case reflect.Float32, reflect.Float64:
		var res floatStr
		if err := json.Unmarshal(data, &res); err != nil {
			return errors.Wrap(err, "while unmarshalling as floatStr")
		}

		field.SetFloat(res.unwrap())
		return nil

	case reflect.Bool:
		if info.hex {
			var res hexutil.Uint
			if err := json.Unmarshal(data, &res); err != nil {
				return errors.Wrap(err, "while unmarshalling as hexutil.Uint")
			}

			field.SetBool(uint64(res) != 0)
			return nil
		}

		if info.num {
			var res string
			if err := json.Unmarshal(data, &res); err != nil {
				return errors.Wrap(err, "while unmarshalling as string")
			}

			field.SetBool(res != "0")
			return nil
		}

		return json.Unmarshal(data, field.Addr().Interface())

	case reflect.Slice:
		return unmarshalSliceRsp(data, field)

	default:
		if err := json.Unmarshal(data, field.Addr().Interface()); err != nil {
			return errors.Wrap(err, "while unmarshalling as value")
		}

		return nil
	}
}
