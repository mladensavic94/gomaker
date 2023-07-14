package gomaker

import (
	"fmt"
	"math/rand"
	"reflect"
	"unsafe"
)

var tag = "gomaker"

type option string

const (
	random option = "rand"
	regex  option = "regex"
	val    option = "val"
)

type Maker struct {
}

func New() *Maker {
	return &Maker{}
}

func (m Maker) Fill(model any) error {
	if err := enforcePointer(reflect.TypeOf(model)); err != nil {
		return err
	}
	return m.fillStruct(model)
}

func (m Maker) fillStruct(model any) error {
	valueOf := reflect.Indirect(reflect.ValueOf(model))
	typeOf := valueOf.Type()
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		if !field.IsExported() {
			continue
		}
		tagValue, ok := field.Tag.Lookup(tag)
		if ok {
			kind := field.Type.Kind()
			switch option(tagValue) {
			case random:
				if err := handleRandom(valueOf.FieldByName(field.Name), kind); err != nil {
					return err
				}
			case regex:
			case val:
			default:
				return fmt.Errorf("option not available %s", tagValue)
			}
		}
	}
	return nil
}

func enforcePointer(model reflect.Type) error {
	if model.Kind() != reflect.Pointer {
		return fmt.Errorf("non-pointer argument")
	}
	return nil
}

func handleRandom(field reflect.Value, kind reflect.Kind) error {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.SetInt(rand.Int63())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		field.SetUint(rand.Uint64())
	case reflect.Float32, reflect.Float64:
		field.SetFloat(rand.Float64())
	case reflect.String:
		field.SetString(randString(rand.Int63n(20)))
	case reflect.Bool:
		field.SetBool(rand.Float64() < 0.5)
	default:
		return fmt.Errorf("kind not supported: %s", kind.String())
	}
	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits
)

func randString(n int64) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
