package gomaker

import (
	"fmt"
	"reflect"
	"strings"
)

func fillFuncSimple(funcMap map[string]func() any, field reflect.Value, tagValue string) error {
	kind := field.Kind()
	funcName := getFuncName(tagValue)
	fn, found := funcMap[funcName]
	if !found {
		return fmt.Errorf("map missing fn %s", funcName)
	}
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res, ok := fn().(int64)
		if !ok {
			return fmt.Errorf("expected int64 got %v", reflect.TypeOf(fn()).Name())
		}
		field.SetInt(res)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		res, ok := fn().(uint64)
		if !ok {
			return fmt.Errorf("expected uint64 got %v", reflect.TypeOf(fn()).Name())
		}
		field.SetUint(res)
	case reflect.Float32, reflect.Float64:
		res, ok := fn().(float64)
		if !ok {
			return fmt.Errorf("expected float64 got %v", reflect.TypeOf(fn()).Name())
		}
		field.SetFloat(res)
	case reflect.Complex64, reflect.Complex128:
		res, ok := fn().(complex128)
		if !ok {
			return fmt.Errorf("expected complex128 got %v", reflect.TypeOf(fn()).Name())
		}
		field.SetComplex(res)
	case reflect.String:
		res, ok := fn().(string)
		if !ok {
			return fmt.Errorf("expected string got %v", reflect.TypeOf(fn()).Name())
		}
		field.SetString(res)
	case reflect.Bool:
		res, ok := fn().(bool)
		if !ok {
			return fmt.Errorf("expected bool got %v", reflect.TypeOf(fn()).Name())
		}
		field.SetBool(res)
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
