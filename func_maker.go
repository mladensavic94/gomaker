package gomaker

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func fillFuncSimple(funcMap map[string]func() string, field reflect.Value, tagValue string) error {
	kind := field.Kind()
	funcName := getFuncName(tagValue)
	fn, found := funcMap[funcName]
	if !found {
		return fmt.Errorf("map missing fn %s", funcName)
	}
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(fn(), 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		i, err := strconv.ParseUint(fn(), 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(fn(), 64)
		if err != nil {
			return err
		}
		field.SetFloat(i)
	case reflect.Complex64, reflect.Complex128:
		i, err := strconv.ParseComplex(fn(), 128)
		if err != nil {
			return err
		}
		field.SetComplex(i)
	case reflect.String:
		field.SetString(fn())
	case reflect.Bool:
		i, err := strconv.ParseBool(fn())
		if err != nil {
			return err
		}
		field.SetBool(i)
	default:
		return fmt.Errorf("kind not supported: %s", kind.String())
	}
	return nil
}

func getFuncName(value string) string {
	value = strings.TrimPrefix(value, "func[")
	value = strings.TrimSuffix(value, "]")
	return value
}
