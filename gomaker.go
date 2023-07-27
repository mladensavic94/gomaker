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
)

type Maker struct {
	rand *rand.Rand
}

func New(options ...func(maker *Maker)) *Maker {
	m := &Maker{rand: rand.New(rand.NewSource(time.Now().Unix()))}
	for _, opt := range options {
		opt(m)
	}
	return m
}

func WithSeed(seed int64) func(maker *Maker) {
	return func(maker *Maker) {
		maker.rand = rand.New(rand.NewSource(seed))
	}
}

func (m Maker) Fill(model any) error {
	if reflect.TypeOf(model).Kind() != reflect.Pointer {
		return fmt.Errorf("non-pointer argument")
	}
	return m.fillStruct(reflect.Indirect(reflect.ValueOf(model)))
}

func (m Maker) fillStruct(valueOf reflect.Value) error {
	typeOf := valueOf.Type()
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		if !field.IsExported() {
			continue
		}
		tagValue := field.Tag.Get(tag)
		kind := field.Type.Kind()
		if kind == reflect.Struct {
			if err := m.fillStruct(valueOf.FieldByName(field.Name)); err != nil {
				return err
			}
		} else {
			if err := m.fillSimple(tagValue, valueOf.FieldByName(field.Name)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m Maker) fillSimple(tagValue string, field reflect.Value) error {
	switch optionValueOf(tagValue) {
	case random:
		if err := fillRandomSimple(m.rand, field, tagValue); err != nil {
			return err
		}
	case regex:
	default:
		return fmt.Errorf("option not available %s", tagValue)
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
	return ""
}
