package gomaker

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"unsafe"
)

type constraints struct {
	min, max int64
	step     float64
}

var defaultConstraints = constraints{min: 1, max: 10, step: 1}
var randomPattern = regexp.MustCompile(`\[(\d*)?;(\d*)?;(\d*\.?\d)?]`)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func (c constraints) Validate(kind reflect.Kind) error {
	if c.min > c.max {
		return errors.New("min bigger then max")
	}
	if c.step < 0 {
		return errors.New("negative step")
	}
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if math.Mod(c.step, 1) != 0 {
			return errors.New("step not whole number for int type")
		}
	}
	return nil
}

func fillRandomSimple(seed int64, field reflect.Value, tagValue string) error {
	c := getOptions(tagValue)
	kind := field.Kind()
	if err := c.Validate(kind); err != nil {
		return err
	}
	r := rand.New(rand.NewSource(seed))
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.SetInt(randInt64(r, c))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		field.SetUint(uint64(randInt64(r, c)))
	case reflect.Float32, reflect.Float64:
		field.SetFloat(randFloat64(r, c))
	case reflect.Complex64, reflect.Complex128:
		field.SetComplex(complex(randFloat64(r, c), randFloat64(r, c)))
	case reflect.String:
		field.SetString(randString(r, randInt64(r, c)))
	case reflect.Bool:
		field.SetBool(r.Float64() < 0.5)
	default:
		return fmt.Errorf("kind not supported: %s", kind.String())
	}
	return nil
}

func getOptions(value string) constraints {
	matches := randomPattern.FindStringSubmatch(value)
	c := defaultConstraints
	if len(matches) == 4 {
		if tmp, err := strconv.Atoi(matches[1]); err == nil {
			c.min = int64(tmp)
		}
		if tmp, err := strconv.Atoi(matches[2]); err == nil {
			c.max = int64(tmp)
		}
		if tmp, err := strconv.ParseFloat(matches[3], 32); err == nil {
			c.step = tmp
		}
	}
	return c
}

func randInt64(r *rand.Rand, in constraints) int64 {
	res := randFloat64(r, in)
	return int64(res)
}

func randFloat64(r *rand.Rand, in constraints) float64 {
	scale := r.Float64()*float64(in.max-in.min) + float64(in.min)
	mod := math.Mod(scale, in.step)
	res := (scale - mod) * in.step
	return res
}

func randString(r *rand.Rand, n int64) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, r.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = r.Int63(), letterIdxMax
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
