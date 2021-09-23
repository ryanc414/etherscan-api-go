package etherscan

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func marshalRequest(req interface{}) (map[string]string, error) {
	reqType := reflect.TypeOf(req)

	if reqType.Kind() != reflect.Struct {
		return nil, errors.New("request must be a struct type")
	}

	res := make(map[string]string)
	reqVal := reflect.ValueOf(req)

	for i := 0; i < reqType.NumField(); i++ {
		key := keyName(reqType.Field(i))
		res[key] = formatValue(reqVal.Field(i))
	}

	return res, nil
}

func keyName(fieldType reflect.StructField) string {
	if tag := fieldType.Tag.Get("etherscan"); tag != "" {
		return tag
	}

	return strings.ToLower(fieldType.Name)
}

func formatValue(fieldVal reflect.Value) string {
	switch v := fieldVal.Interface().(type) {
	default:
		return fmt.Sprint(v)
	}
}
