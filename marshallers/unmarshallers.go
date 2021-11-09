package marshallers

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

type BigInt big.Int

func (b *BigInt) UnmarshalJSON(data []byte) error {
	var intStr string
	if err := json.Unmarshal(data, &intStr); err != nil {
		return err
	}

	val, ok := new(big.Int).SetString(intStr, 10)
	if !ok {
		return errors.Errorf("cannot parse %s as big.Int", intStr)
	}

	*b = BigInt(*val)
	return nil
}

func (b *BigInt) Unwrap() *big.Int {
	return (*big.Int)(b)
}

type commaDecimal decimal.Decimal

func (d *commaDecimal) UnmarshalJSON(data []byte) error {
	var decStr string
	if err := json.Unmarshal(data, &decStr); err != nil {
		return err
	}

	decStr = strings.Replace(decStr, ",", "", -1)

	dec, err := decimal.NewFromString(decStr)
	if err != nil {
		return err
	}

	*d = commaDecimal(dec)
	return nil
}

func (d commaDecimal) unwrap() decimal.Decimal {
	return decimal.Decimal(d)
}

type UintStr uint64

func (u *UintStr) UnmarshalJSON(data []byte) error {
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

	*u = UintStr(val)
	return nil
}

func (u UintStr) Unwrap() uint64 {
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

type dateTimestamp time.Time

func (t *dateTimestamp) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}

	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return err
	}

	*t = dateTimestamp(date)
	return nil
}

func (t dateTimestamp) unwrap() time.Time {
	return time.Time(t)
}

func UnmarshalResponse(data []byte, v interface{}) error {
	rspType := reflect.TypeOf(v)
	if rspType.Kind() != reflect.Ptr {
		return errors.New("value must be a pointer")
	}

	var u json.Unmarshaler
	if rspType.Implements(reflect.TypeOf(&u).Elem()) {
		return json.Unmarshal(data, v)
	}

	rspVal := reflect.ValueOf(v).Elem()

	switch rspVal.Kind() {
	case reflect.Struct:
		return unmarshalStructRsp(data, rspVal)

	case reflect.Slice:
		if rspVal.Type().Elem().Kind() != reflect.Struct {
			return errors.New("only slices of structs are allowed")
		}

		return unmarshalSliceRsp(data, rspVal, nil)

	default:
		return json.Unmarshal(data, v)
	}
}

func unmarshalSliceRsp(data []byte, v reflect.Value, info *tagInfo) error {
	if info != nil && info.sep {
		return unmarshalStringSepSlice(data, v, ",", info)
	}

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

func unmarshalStringSepSlice(data []byte, v reflect.Value, sep string, info *tagInfo) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return errors.Wrap(err, "while unmarshalling as string")
	}

	substrs := strings.Split(str, sep)
	slice := reflect.MakeSlice(v.Type(), len(substrs), len(substrs))

	for i := range substrs {
		el := slice.Index(i)

		val := fmt.Sprintf("\"%s\"", substrs[i])
		if err := setFieldValue(el, []byte(val), info); err != nil {
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
			log.Warn().Msgf("no field with name %s in response data", name)
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

func unmarshalField(
	data []byte,
	into interface{},
	field reflect.Value,
	setter func(interface{}),
) error {
	if err := json.Unmarshal(data, into); err != nil {
		return errors.Wrap(err, "while unmarshalling json value")
	}

	if setter == nil {
		return nil
	}

	setter(into)
	return nil
}

func setFieldValue(field reflect.Value, data []byte, info *tagInfo) error {
	if string(data) == "\"\"" || string(data) == "null" {
		return nil
	}

	if _, ok := field.Interface().([]byte); !ok && field.Kind() == reflect.Slice {
		return unmarshalSliceRsp(data, field, info)
	}

	into, setter := getTypeUnmarshler(field, data, info)
	if into == nil {
		return nil
	}

	return unmarshalField(data, into, field, setter)
}

func setDirect(v interface{}, field reflect.Value) {
	field.Set(reflect.ValueOf(v))
}

func getTypeUnmarshler(
	field reflect.Value, data []byte, info *tagInfo,
) (interface{}, func(interface{})) {
	iField := field.Interface()
	if _, ok := iField.(*big.Int); ok {
		if info.hex {
			return new(hexutil.Big), func(v interface{}) {
				setDirect(v.(*hexutil.Big).ToInt(), field)
			}
		}

		if info.num {
			return new(big.Int), func(v interface{}) {
				setDirect(v.(*big.Int), field)
			}
		}

		return new(BigInt), func(v interface{}) {
			setDirect(v.(*BigInt).Unwrap(), field)
		}
	}

	if _, ok := iField.(decimal.Decimal); ok && info.comma {
		return new(commaDecimal), func(v interface{}) {
			setDirect(v.(*commaDecimal).unwrap(), field)
		}
	}

	if _, ok := iField.(time.Time); ok {
		if info.hex {
			return new(hexTimestamp), func(v interface{}) {
				setDirect(v.(*hexTimestamp).unwrap(), field)
			}
		}

		if info.date {
			return new(dateTimestamp), func(v interface{}) {
				setDirect(v.(*dateTimestamp).unwrap(), field)
			}
		}

		return new(unixTimestamp), func(v interface{}) {
			setDirect(v.(*unixTimestamp).unwrap(), field)
		}
	}

	if _, ok := iField.([]byte); ok {
		if string(data) == "\"deprecated\"" {
			return nil, nil
		}

		return new(hexutil.Bytes), func(v interface{}) {
			setDirect([]byte(*v.(*hexutil.Bytes)), field)
		}
	}

	switch field.Kind() {
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		if info.num {
			return field.Addr().Interface(), nil
		}

		if info.hex {
			return new(hexUint), func(v interface{}) {
				field.SetUint(uint64(*v.(*hexUint)))
			}
		}

		return new(UintStr), func(v interface{}) {
			field.SetUint(uint64(*v.(*UintStr)))
		}

	case reflect.Bool:
		if info.hex {
			return new(hexutil.Uint), func(v interface{}) {
				val := uint64(*v.(*hexutil.Uint))
				field.SetBool(val != 0)
			}
		}

		if info.num {
			return new(UintStr), func(v interface{}) {
				val := uint64(*v.(*UintStr))
				field.SetBool(val != 0)
			}
		}

		if info.str {
			return new(string), func(v interface{}) {
				val := *v.(*string)
				field.SetBool(val == "true")
			}
		}

		return field.Addr().Interface(), nil

	default:
		return field.Addr().Interface(), nil
	}
}
