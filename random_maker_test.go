package gomaker

import "testing"

func Test_getOptions(t *testing.T) {
	tests := []struct {
		name  string
		arg   string
		want  int64
		want1 int64
		want2 float64
	}{
		{
			"nothing",
			"rand",
			1,
			10,
			1,
		},
		{
			"full",
			"rand[2;11;.5]",
			2,
			11,
			0.5,
		},
		{
			"min only",
			"rand[2;;]",
			2,
			10,
			1,
		},
		{
			"max only",
			"rand[;11;]",
			1,
			11,
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getOptions(tt.arg)
			if got.min != tt.want {
				t.Errorf("getOptions() got = %v, want %v", got.min, tt.want)
			}
			if got.max != tt.want1 {
				t.Errorf("getOptions() got1 = %v, want %v", got.max, tt.want1)
			}
			if got.step != tt.want2 {
				t.Errorf("getOptions() got2 = %v, want %v", got.step, tt.want2)
			}
		})
	}
}

func Test_randInt64(t *testing.T) {
	tests := []struct {
		name string
		args constraints
		want func(in int64) bool
	}{
		{
			"normal",
			constraints{min: 1, max: 10, step: 1},
			func(in int64) bool {
				return in >= 1 && in < 10 && in%1 == 0
			},
		},
		{
			"wrong step",
			constraints{min: 1, max: 10, step: 0.1},
			func(in int64) bool {
				return in == 0
			},
		},
		{
			"same length",
			constraints{min: 10, max: 10, step: 1},
			func(in int64) bool {
				return in == 10
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := randInt64(tt.args)
			if !tt.want(got) {
				t.Errorf("randInt64() = %v", got)
			}
		})
	}
}
