package gomaker_test

import (
	"errors"
	"fmt"
	"gomaker"
	"testing"
	"time"
)

func TestMaker_TDD_random(t *testing.T) {
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
		DummyId int64 `gomaker:"test123"`
	}

	maker := gomaker.New()
	tests := []struct {
		name   string
		arg    any
		err    error
		sanity func(in *dummy) error
	}{
		{
			"pass non pointer",
			dummy{},
			fmt.Errorf("non-pointer argument"),
			nil,
		},
		{
			"pass unknown",
			&unknown{},
			fmt.Errorf("option not available test123"),
			nil,
		},
		{
			"happy path",
			&dummy{},
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
			err := maker.Fill(tt.arg)
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

func Test_repeatable_generation(t *testing.T) {
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

func TestMaker_TDD_regex(t *testing.T) {
	type dummy struct {
		DummyString string `gomaker:"regex[^\\d+$]"`
	}
	maker := gomaker.New()
	tests := []struct {
		name   string
		arg    any
		err    error
		sanity func(in *dummy) error
	}{
		{
			"pass non pointer",
			dummy{},
			errors.New("non-pointer argument"),
			nil,
		},
		{
			"happy path",
			&dummy{},
			nil,
			func(in *dummy) error {
				if in.DummyString == "" {
					return errors.New("string not assigned")
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

func TestMaker_TDD_race(t *testing.T) {
	type dummy struct {
		DummyString string `gomaker:"regex[^\\d+$]"`
	}
	maker := gomaker.New()
	go maker.Fill(&dummy{})
	go maker.Fill(&dummy{})
	go maker.Fill(&dummy{})

	time.Sleep(time.Second)

}
