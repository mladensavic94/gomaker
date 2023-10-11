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
	funcMap map[string]func() any
	fields  map[string]any
}

func New(options ...func(maker *Maker)) *Maker {
	m := &Maker{seed: time.Now().Unix(), funcMap: map[string]func() any{}}
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

func WithFuncMap(f map[string]func() any) func(maker *Maker) {
	return func(maker *Maker) {
		maker.funcMap = f
	}
}

func WithFieldsMapping(f map[string]any) func(maker *Maker) {
	return func(maker *Maker) {
		maker.fields = f
	}
}

func (m *Maker) Fill(model any) error {
	if reflect.TypeOf(model).Kind() != reflect.Pointer {
		return fmt.Errorf("non-pointer argument")
	}
	value := reflect.Indirect(reflect.ValueOf(model))
	if len(m.fields) == 0 {
		res, err := buildGraph(value)
		if err != nil {
			return err
		}
		m.fields = res
	}
	return m.fillStruct(rand.New(rand.NewSource(m.seed)), value, m.fields)
}

func buildGraph(valueOf reflect.Value) (map[string]any, error) {
	graph := make(map[string]any)
	var err error
	typeOf := valueOf.Type()
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		if !field.IsExported() {
			continue
		}
		tagValue := field.Tag.Get(tag)
		kind := field.Type.Kind()
		if kind == reflect.Struct {
			m, _ := buildGraph(valueOf.FieldByName(field.Name))
			graph[field.Name] = m
		} else if kind == reflect.Slice {
			if valueOf.Field(i).Type().Elem().Kind() == reflect.Struct {
				m, _ := buildGraph(valueOf.Field(i).Index(0))
				graph[field.Name] = m
			} else {
				if tagValue == "" {
					continue
				}
				graph[field.Name] = tagValue
			}

		} else {
			if tagValue == "" {
				continue
			}
			graph[field.Name] = tagValue
		}
	}
	return graph, err
}

func (m *Maker) fillStruct(r *rand.Rand, valueOf reflect.Value, fields map[string]any) error {
	var err error
	for key, val := range fields {
		fieldsKey := reflect.TypeOf(val).Kind()
		if fieldsKey == reflect.String {
			tagValue := val.(string)
			if valueOf.Kind() == reflect.Slice {
				return m.fillSlice(r, tagValue, valueOf, fields)
			}
			fieldByName := valueOf.FieldByName(key)
			kind := fieldByName.Kind()
			if kind == reflect.Struct {
				err = m.fillStruct(r, fieldByName, fields)
			} else if kind == reflect.Slice {
				err = m.fillSlice(r, tagValue, fieldByName, fields)
			} else {
				err = m.fillSimple(r, tagValue, fieldByName)
			}
		} else if fieldsKey == reflect.Map {
			err = m.fillStruct(r, valueOf.FieldByName(key), fields[key].(map[string]any))
		} else {
			return fmt.Errorf("unrecognized type %v", fieldsKey)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Maker) fillSimple(r *rand.Rand, tagValue string, field reflect.Value) error {
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

func (m *Maker) fillSlice(r *rand.Rand, tagValue string, field reflect.Value, fields map[string]any) error {
	for i := 0; i < field.Len(); i++ {
		index := field.Index(i)
		var err error
		if index.Kind() == reflect.Struct {
			err = m.fillStruct(r, index, fields)
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
