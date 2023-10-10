package gomaker_test

import (
	"errors"
	"fmt"
	"gomaker"
	"math/rand"
	"testing"
)

func TestMaker_random_with_tags(t *testing.T) {
	t.Parallel()
	type inner struct {
		InnerInt int32 `gomaker:"rand"`
	}
	type dummy struct {
		DummyId      int64      `gomaker:"rand[1;10;1]"`
		DummyString  string     `gomaker:"rand"`
		DummyComplex complex128 `gomaker:"rand"`
		Inner        inner
	}
	type unknown struct {
		UnknownId int64 `gomaker:"test123"`
	}

	tests := []struct {
		name   string
		arg    any
		maker  *gomaker.Maker
		err    error
		sanity func(in *dummy) error
	}{
		{
			"pass non pointer",
			dummy{},
			gomaker.New(),
			fmt.Errorf("non-pointer argument"),
			nil,
		},
		{
			"pass unknown",
			&unknown{},
			gomaker.New(),
			fmt.Errorf("option not available test123"),
			nil,
		},
		{
			"happy path",
			&dummy{},
			gomaker.New(),
			nil,
			func(in *dummy) error {
				if in.DummyId == 0 {
					return errors.New("int not assigned")
				}
				if in.DummyString == "" {
					return errors.New("string not assigned")
				}
				if in.DummyComplex == 0 {
					return errors.New("complex not assigned")
				}
				if in.Inner.InnerInt == 0 {
					return errors.New("inner not assigned")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.maker.Fill(tt.arg)
			if err != nil {
				if tt.err != nil && err.Error() != tt.err.Error() {
					t.Fatalf("expected: %v, got: %v", tt.err, err)
				}
			}
			if tt.sanity != nil {
				err = tt.sanity(tt.arg.(*dummy))
				if err != nil {
					t.Fatalf("sanity check failed: %v", err)
				}
			}
		})
	}
}

func TestMaker_random(t *testing.T) {
	t.Parallel()
	type inner struct {
		InnerInt int32
	}
	type dummy struct {
		DummyId      int64
		DummyString  string
		DummyComplex complex128
		Inner        inner
	}
	type unknown struct {
		UnknownId int64
	}

	tests := []struct {
		name   string
		arg    any
		maker  *gomaker.Maker
		err    error
		sanity func(in *dummy) error
	}{
		{
			"pass non pointer",
			dummy{},
			gomaker.New(),
			fmt.Errorf("non-pointer argument"),
			nil,
		},
		{
			"pass unknown",
			&unknown{},
			gomaker.New(gomaker.WithFieldsMapping(map[string]any{"UnknownId": "test123"})),
			fmt.Errorf("option not available test123"),
			nil,
		},
		{
			"happy path",
			&dummy{},
			gomaker.New(gomaker.WithFieldsMapping(map[string]any{"DummyId": "rand[1;10;1]", "DummyString": "rand", "DummyComplex": "rand", "Inner": map[string]any{"InnerInt": "rand"}})),
			nil,
			func(in *dummy) error {
				if in.DummyId == 0 {
					return errors.New("int not assigned")
				}
				if in.DummyString == "" {
					return errors.New("string not assigned")
				}
				if in.DummyComplex == 0 {
					return errors.New("complex not assigned")
				}
				if in.Inner.InnerInt == 0 {
					return errors.New("inner not assigned")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.maker.Fill(tt.arg)
			if err != nil {
				if tt.err != nil && err.Error() != tt.err.Error() {
					t.Fatalf("expected: %v, got: %v", tt.err, err)
				}
			}
			if tt.sanity != nil {
				err = tt.sanity(tt.arg.(*dummy))
				if err != nil {
					t.Fatalf("sanity check failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkRandFill(b *testing.B) {
	type dummy struct {
		DummyId     int64  `gomaker:"rand[1;100;5]"`
		DummyString string `gomaker:"rand[10;10;]"`
	}

	maker := gomaker.New()
	model := &dummy{}
	for i := 0; i < b.N; i++ {
		err := maker.Fill(model)
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkRandFill_WithPreloadMapping(b *testing.B) {
	type dummy struct {
		DummyId     int64
		DummyString string
	}

	maker := gomaker.New(gomaker.WithFieldsMapping(map[string]any{"DummyId": "rand[1;100;5]", "DummyString": "rand[10;10;]"}))
	model := &dummy{}
	for i := 0; i < b.N; i++ {
		err := maker.Fill(model)
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkRegexFill_WithPreloadMapping(b *testing.B) {
	type dummy struct {
		DummyId     int64
		DummyString string
	}

	maker := gomaker.New(gomaker.WithFieldsMapping(map[string]any{"DummyId": `regex[[0-9]{10}]`, "DummyString": `regex[(abc)+]`}))
	model := &dummy{}
	for i := 0; i < b.N; i++ {
		err := maker.Fill(model)
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkRegexFill(b *testing.B) {
	type dummy struct {
		DummyId     int64  `gomaker:"regex[[0-9]{10}]"`
		DummyString string `gomaker:"regex[(abc)+]"`
	}

	maker := gomaker.New()
	model := &dummy{}
	for i := 0; i < b.N; i++ {
		err := maker.Fill(model)
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkFuncFill(b *testing.B) {
	type dummy struct {
		DummyId     int64  `gomaker:"func[randomInt]"`
		DummyString string `gomaker:"func[flatStr]"`
	}

	maker := gomaker.New(gomaker.WithFuncMap(map[string]func() string{"randomInt": func() string {
		return fmt.Sprint(rand.Int())
	}, "flatStr": func() string {
		return "123456"
	}}))
	model := &dummy{}
	for i := 0; i < b.N; i++ {
		err := maker.Fill(model)
		if err != nil {
			println(err.Error())
			b.FailNow()
		}
	}
}

func BenchmarkFuncFill_WithPreloadMapping(b *testing.B) {
	type dummy struct {
		DummyId     int64
		DummyString string
	}

	maker := gomaker.New(gomaker.WithFuncMap(map[string]func() string{"randomInt": func() string {
		return fmt.Sprint(rand.Int())
	}, "flatStr": func() string {
		return "123456"
	}}), gomaker.WithFieldsMapping(map[string]any{"DummyId": "func[randomInt]", "DummyString": "func[flatStr]"}))
	model := &dummy{}
	for i := 0; i < b.N; i++ {
		err := maker.Fill(model)
		if err != nil {
			println(err.Error())
			b.FailNow()
		}
	}
}

func Test_repeatable_generation(t *testing.T) {
	t.Parallel()
	type dummy struct {
		DummyId     int64  `gomaker:"rand[1;100;5]"`
		DummyString string `gomaker:"rand[10;10;]"`
	}
	m1 := gomaker.New(gomaker.WithSeed(123))
	model1 := &dummy{}
	if err := m1.Fill(model1); err != nil {
		t.Errorf(err.Error())
	}
	m2 := gomaker.New(gomaker.WithSeed(123))
	model2 := &dummy{}
	if err := m2.Fill(model2); err != nil {
		t.Errorf(err.Error())
	}

	if model1.DummyId != model2.DummyId {
		t.Errorf("got %v expected %v", model1.DummyId, model2.DummyId)
	}
	if model1.DummyString != model2.DummyString {
		t.Errorf("got %v expected %v", model1.DummyString, model2.DummyString)
	}
}

func TestMaker_regex(t *testing.T) {
	t.Parallel()
	type dummy struct {
		DummyString string
		DummyInt    int32
	}
	type failRegex struct {
		Str string
	}
	tests := []struct {
		name   string
		arg    any
		maker  *gomaker.Maker
		err    error
		sanity func(in *dummy) error
	}{
		{
			"pass non pointer",
			dummy{},
			gomaker.New(),
			errors.New("non-pointer argument"),
			nil,
		},
		{
			"happy path",
			&dummy{},
			gomaker.New(gomaker.WithFieldsMapping(map[string]any{"DummyString": `regex[^\\d+$]`, "DummyInt": `regex[[0-9]{10}]`})),
			nil,
			func(in *dummy) error {
				if in.DummyString == "" {
					return errors.New("string not assigned")
				}
				if in.DummyInt == 0 {
					return errors.New("int not assigned")
				}
				return nil
			},
		},
		{
			"fail regex",
			&failRegex{},
			gomaker.New(gomaker.WithFieldsMapping(map[string]any{"Str": `regexalmost[]`})),
			errors.New("regex validation failed"),
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.maker.Fill(tt.arg)
			if err != nil {
				if (tt.err != nil && err.Error() != tt.err.Error()) || tt.err == nil {
					t.Fatalf("expected: %v, got: %v", tt.err, err)
				}
			}
			if tt.sanity != nil {
				err = tt.sanity(tt.arg.(*dummy))
				if err != nil {
					t.Fatalf("sanity check failed: %v", err)
				}
			}
		})
	}
}

func TestMaker_func(t *testing.T) {
	t.Parallel()
	type dummy struct {
		DummyString  string
		DummyComplex complex128
		DummyId      int64
		DummyFloat   float32
		DummyBool    bool
		DummyUint    uint64
	}
	type unknown struct {
		DummyId  int64 `gomaker:"func[missing]"`
		DummyInt int64
	}
	funcMap := map[string]func() string{"test": func() string {
		return "123"
	}, "complex": func() string {
		return "1+1i"
	}, "int64": func() string {
		return "12"
	}, "float": func() string {
		return "1.2"
	}, "bool": func() string {
		return "t"
	}}
	tests := []struct {
		name   string
		arg    any
		maker  *gomaker.Maker
		err    error
		sanity func(in *dummy) error
	}{
		{
			"missing func",
			&unknown{},
			gomaker.New(gomaker.WithFuncMap(funcMap)),
			errors.New("map missing fn missing"),
			nil,
		},
		{
			"happy path",
			&dummy{},
			gomaker.New(gomaker.WithFuncMap(funcMap), gomaker.WithFieldsMapping(map[string]any{"DummyUint": "func[int64]", "DummyBool": "func[bool]", "DummyFloat": "func[float]", "DummyId": "func[int64]", "DummyComplex": "func[complex]", "DummyString": "func[test]"})),
			nil,
			func(in *dummy) error {
				if in.DummyString != "123" {
					return errors.New("string not assigned")
				}
				if in.DummyComplex != complex(1, 1) {
					return errors.New("complex not assigned")
				}
				if in.DummyId != 12 {
					return errors.New("int not assigned")
				}
				if in.DummyFloat != 1.2 {
					return errors.New("float not assigned")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.maker.Fill(tt.arg)
			if err != nil {
				if (tt.err != nil && err.Error() != tt.err.Error()) || tt.err == nil {
					t.Fatalf("expected: %v, got: %v", tt.err, err)
				}
			}
			if tt.sanity != nil {
				err = tt.sanity(tt.arg.(*dummy))
				if err != nil {
					t.Fatalf("sanity check failed: %v", err)
				}
			}
		})
	}
}

func TestMaker_customTypes(t *testing.T) {
	t.Parallel()
	type wrapper string
	type dummy struct {
		WrapRand  wrapper
		WrapRegex wrapper
		WrapFunc  wrapper
	}
	funcMap := map[string]func() string{"test": func() string {
		return "123"
	}}
	tests := []struct {
		name   string
		arg    any
		maker  *gomaker.Maker
		err    error
		sanity func(in *dummy) error
	}{
		{
			"customTypes",
			&dummy{},
			gomaker.New(gomaker.WithFuncMap(funcMap), gomaker.WithFieldsMapping(map[string]any{"WrapRand": "rand[10;10;]", "WrapRegex": `regex[^test$]`, "WrapFunc": "func[test]"})),
			nil,
			func(in *dummy) error {
				if in.WrapRand == "" {
					return errors.New("rand not assigned")
				}
				if in.WrapRegex == "" {
					return errors.New("regex got assigned")
				}
				if in.WrapFunc == "" {
					return errors.New("func not assigned")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.maker.Fill(tt.arg)
			if err != nil {
				if (tt.err != nil && err.Error() != tt.err.Error()) || tt.err == nil {
					t.Fatalf("expected: %v, got: %v", tt.err, err)
				}
			}
			if tt.sanity != nil {
				err = tt.sanity(tt.arg.(*dummy))
				if err != nil {
					t.Fatalf("sanity check failed: %v", err)
				}
			}
		})
	}
}

func TestMaker_errors(t *testing.T) {
	t.Parallel()
	type dummy struct {
		DummyId int64 `gomaker:"func[bool]"`
	}
	maker := gomaker.New(gomaker.WithFuncMap(map[string]func() string{"bool": func() string {
		return "t"
	}}))
	tests := []struct {
		name   string
		arg    any
		err    error
		sanity func(in *dummy) error
	}{
		{
			"errors",
			&dummy{},
			errors.New("strconv.ParseInt: parsing \"t\": invalid syntax"),
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := maker.Fill(tt.arg)
			if err != nil {
				if (tt.err != nil && err.Error() != tt.err.Error()) || tt.err == nil {
					t.Fatalf("expected: %v, got: %v", tt.err, err)
				}
			}
			if tt.sanity != nil {
				err = tt.sanity(tt.arg.(*dummy))
				if err != nil {
					t.Fatalf("sanity check failed: %v", err)
				}
			}
		})
	}
}

func TestMaker_slices(t *testing.T) {
	t.Parallel()
	type inner struct {
		Floats float64 `gomaker:"rand[1;100;0.1]"`
	}
	type dummy struct {
		Ints   []int64  `gomaker:"func[int64]"`
		Strs   []string `gomaker:"regex[.{5}]"`
		Inners []inner
	}
	maker := gomaker.New(gomaker.WithFuncMap(map[string]func() string{"int64": func() string {
		return "1"
	}}))
	tests := []struct {
		name   string
		arg    any
		err    error
		sanity func(in *dummy) error
	}{
		{
			"slices",
			&dummy{Ints: make([]int64, 10), Strs: make([]string, 5), Inners: make([]inner, 2)},
			nil,
			func(in *dummy) error {
				if len(in.Ints) == 0 && in.Ints[0] != 0 {
					return errors.New("ints not assigned")
				}
				if len(in.Strs) == 0 && in.Strs[0] != "" {
					return errors.New("strs not assigned")
				}
				if len(in.Inners) == 0 && in.Inners[0].Floats != 0 {
					return errors.New("inners not assigned")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := maker.Fill(tt.arg)
			if err != nil {
				if (tt.err != nil && err.Error() != tt.err.Error()) || tt.err == nil {
					t.Fatalf("expected: %v, got: %v", tt.err, err)
				}
			}
			if tt.sanity != nil {
				err = tt.sanity(tt.arg.(*dummy))
				if err != nil {
					t.Fatalf("sanity check failed: %v", err)
				}
			}
		})
	}
}
