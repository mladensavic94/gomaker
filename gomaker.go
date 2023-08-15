package gomaker

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

var tag = "gomaker"

type option string

const (
	random option = "rand"
	regex  option = "regex"
	rel    option = "rel"
	fc     option = "func"
)

type Maker struct {
	seed    int64
	funcMap map[string]func() string
}

func New(options ...func(maker *Maker)) *Maker {
	m := &Maker{seed: time.Now().Unix(), funcMap: map[string]func() string{}}
	for _, opt := range options {
		opt(m)
	}
	return m
}

func WithSeed(seed int64) func(maker *Maker) {
	return func(maker *Maker) {
		maker.seed = seed
	}
}

func WithFuncMap(f map[string]func() string) func(maker *Maker) {
	return func(maker *Maker) {
		maker.funcMap = f
	}
}

func (m Maker) Fill(model any) error {
	if reflect.TypeOf(model).Kind() != reflect.Pointer {
		return fmt.Errorf("non-pointer argument")
	}
	return m.fillStruct(nil, reflect.Indirect(reflect.ValueOf(model)))
}

func (m Maker) fillStruct(r *rand.Rand, valueOf reflect.Value) error {
	typeOf := valueOf.Type()
	if r == nil {
		r = rand.New(rand.NewSource(m.seed))
	}
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		if !field.IsExported() {
			continue
		}
		tagValue := field.Tag.Get(tag)
		kind := field.Type.Kind()
		var err error
		if kind == reflect.Struct {
			err = m.fillStruct(r, valueOf.FieldByName(field.Name))
		} else if kind == reflect.Slice {
			err = m.fillSlice(r, tagValue, valueOf.FieldByName(field.Name))
		} else {
			err = m.fillSimple(r, tagValue, valueOf.FieldByName(field.Name))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (m Maker) fillSimple(r *rand.Rand, tagValue string, field reflect.Value) error {
	switch optionValueOf(tagValue) {
	case random:
		if err := fillRandomSimple(r, field, tagValue); err != nil {
			return err
		}
	case regex:
		if err := fillRegexSimple(r, field, tagValue); err != nil {
			return err
		}
	case fc:
		if err := fillFuncSimple(m.funcMap, field, tagValue); err != nil {
			return err
		}
	default:
		return fmt.Errorf("option not available %s", tagValue)
	}
	return nil
}

func (m Maker) fillSlice(r *rand.Rand, tagValue string, field reflect.Value) error {
	for i := 0; i < field.Len(); i++ {
		index := field.Index(i)
		var err error
		if index.Kind() == reflect.Struct {
			err = m.fillStruct(r, index)
		} else {
			err = m.fillSimple(r, tagValue, index)
		}
		if err != nil {
			return err
		}

	}
	return nil
}

func optionValueOf(in string) option {
	if strings.HasPrefix(in, string(random)) {
		return random
	}
	if strings.HasPrefix(in, string(regex)) {
		return regex
	}
	if strings.HasPrefix(in, string(fc)) {
		return fc
	}
	return ""
}
