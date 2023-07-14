package gomaker_test

import (
	"errors"
	"fmt"
	"gomaker"
	"testing"
)

func TestMaker_TDD(t *testing.T) {
	type dummy struct {
		DummyId     int64  `gomaker:"rand"`
		DummyString string `gomaker:"rand"`
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
			"first",
			&dummy{},
			nil,
			func(in *dummy) error {
				if in.DummyId == 0 {
					return errors.New("int not assigned")
				}
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
		DummyId     int64  `gomaker:"rand"`
		DummyString string `gomaker:"rand"`
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
