package gomaker

import (
	"math/rand"
	"regexp"
	"regexp/syntax"
	"testing"
)

func Test_generate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		regex string
	}{
		{
			"alphanumerics",
			"^[a-zA-Z0-9]*$",
		},
		{
			"5 dots",
			`\.{5}`,
		},
		{
			"any char",
			`.{15}`,
		},
		{
			"repeater",
			`^([A-Z]\.[a-z])*$`,
		},
		{
			"charClass length",
			`[-a-zA-Z0-9@:%._\+~#=]{2,2}`,
		},
		{
			"emails",
			`^([a-zA-Z0-9._%-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6})*$`,
		},
		{
			"time",
			"^(0?[1-9]|1[0-2]):[0-5][0-9]",
		},
		{
			"url",
			`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}([-a-zA-Z0-9@:%_\+.~#?&\/\/=])`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedRegex, err := syntax.Parse(tt.regex, syntax.Perl)
			if err != nil {
				t.Errorf("%v got: %v", tt.name, err)
				return
			}
			r := rand.New(rand.NewSource(12345))
			s, err := generate(r, parsedRegex)
			if err != nil {
				t.Errorf("%v got: %v", tt.name, err)
				return
			}
			m, err := regexp.MatchString(tt.regex, s)
			if !m || err != nil {
				t.Errorf("%v got: %v %v %v", tt.name, err, s, tt.regex)
				return
			}
		})
	}
}
