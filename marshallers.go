package etherscan

import (
	"fmt"
	"log"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func marshalRequest(req interface{}) map[string]string {
	reqType := reflect.TypeOf(req)
	reqVal := reflect.ValueOf(req)

	if reqType.Kind() == reflect.Ptr {
		reqType = reqType.Elem()
		reqVal = reqVal.Elem()
	}

	res := make(map[string]string)

	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)
		info := parseTag(field)
		key := keyName(field, &info)
		res[key] = formatValue(reqVal.Field(i), &info)
	}

	return res
}

func keyName(fieldType reflect.StructField, info *tagInfo) string {
	if info.name != "" {
		return info.name
	}

	return strings.ToLower(fieldType.Name)
}

func formatValue(fieldVal reflect.Value, info *tagInfo) string {
	iVal := fieldVal.Interface()
	if v, ok := iVal.([]byte); ok {
		return hexutil.Encode(v)
	}

	if v, ok := iVal.(*big.Int); ok && info.hex {
		return hexutil.EncodeBig(v)
	}

	switch fieldVal.Kind() {
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		if info.hex {
			return hexutil.EncodeUint64(fieldVal.Uint())
		}

		return strconv.FormatUint(fieldVal.Uint(), 10)

	default:
		return fmt.Sprint(fieldVal.Interface())
	}
}

type tagInfo struct {
	name string
	hex  bool
}

func parseTag(fieldType reflect.StructField) tagInfo {
	rawTag := fieldType.Tag.Get("etherscan")
	log.Print("rawTag: ", rawTag)
	items := strings.Split(rawTag, ",")

	var info tagInfo
	if len(items) == 0 {
		return info
	}

	info.name = items[0]

	for i := 1; i < len(items); i++ {
		if items[i] == "hex" {
			info.hex = true
		}
	}

	return info
}
