package gomaker

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"regexp/syntax"
	"strconv"
	"strings"
)

var regexPattern = regexp.MustCompile(`^regex\[.*]$`)

var generationFailed = errors.New("generator failed")

func fillRegexSimple(r *rand.Rand, field reflect.Value, tagValue string) error {
	if !regexPattern.MatchString(tagValue) {
		return errors.New("regex validation failed")
	}

	parsedRegex, err := getParsedRegex(tagValue)
	if err != nil {
		return errors.New("regex parse failed")
	}

	result, err := generate(r, parsedRegex)
	if err != nil {
		return err
	}

	kind := field.Kind()
	switch kind {
	case reflect.String:
		field.SetString(result)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, _ := strconv.ParseInt(result, 10, 64)
		field.SetInt(i)
	default:
		return fmt.Errorf("kind not supported: %s", kind.String())
	}
	return nil
}

func getParsedRegex(value string) (*syntax.Regexp, error) {
	value = strings.TrimPrefix(value, "regex[")
	value = strings.TrimSuffix(value, "]")
	parse, err := syntax.Parse(value, syntax.Perl)
	if err != nil {
		return nil, err
	}
	return parse.Simplify(), nil
}

func generate(r *rand.Rand, parsedRegex *syntax.Regexp) (string, error) {
	switch parsedRegex.Op {
	case syntax.OpStar:
		return repeatingGenerator(r, parsedRegex.Sub[0], 0, 10)
	case syntax.OpPlus:
		return repeatingGenerator(r, parsedRegex.Sub[0], 1, 10)
	case syntax.OpQuest:
		return repeatingGenerator(r, parsedRegex.Sub[0], 0, 1)
	case syntax.OpRepeat:
		return repeatingGenerator(r, parsedRegex.Sub[0], parsedRegex.Min, parsedRegex.Max)
	case syntax.OpAlternate:
		return alternateGenerator(r, parsedRegex)
	case syntax.OpCharClass:
		return charClassGenerator(r, parsedRegex)
	case syntax.OpCapture:
		return captureGenerator(r, parsedRegex)
	case syntax.OpAnyChar, syntax.OpAnyCharNotNL:
		return charGenerator(r, parsedRegex)
	case syntax.OpLiteral:
		return literalGenerator(r, parsedRegex)
	case syntax.OpConcat:
		return concatGenerator(r, parsedRegex)
	case syntax.OpEndText, syntax.OpEndLine, syntax.OpBeginLine, syntax.OpBeginText, syntax.OpNoWordBoundary, syntax.OpWordBoundary, syntax.OpEmptyMatch:
		return emptyGenerator(r, parsedRegex)
	default:
		return "", fmt.Errorf("op didnt match %s", parsedRegex.Op.String())
	}
}

func repeatingGenerator(r *rand.Rand, parsedRegex *syntax.Regexp, min, max int) (string, error) {
	var buff bytes.Buffer
	repeat := r.Intn(max-min+1) + min
	for i := 0; i < repeat; i++ {
		s, err := generate(r, parsedRegex)
		if err != nil {
			return "", err
		}
		buff.WriteString(s)
	}
	return buff.String(), nil
}

func literalGenerator(_ *rand.Rand, parsedRegex *syntax.Regexp) (string, error) {
	if parsedRegex != nil {
		return string(parsedRegex.Rune), nil
	}
	return "", generationFailed
}

func charGenerator(r *rand.Rand, _ *syntax.Regexp) (string, error) {
	return randString(r, 1), nil
}
func emptyGenerator(_ *rand.Rand, _ *syntax.Regexp) (string, error) {
	return "", nil
}

func concatGenerator(r *rand.Rand, parsedRegex *syntax.Regexp) (string, error) {
	var buff bytes.Buffer
	for _, subRegex := range parsedRegex.Sub {
		s, err := generate(r, subRegex)
		if err != nil {
			return "", generationFailed
		}
		buff.WriteString(s)

	}
	return buff.String(), nil
}

func charClassGenerator(r *rand.Rand, parsedRegex *syntax.Regexp) (string, error) {
	runes := parsedRegex.Rune
	rLen := len(runes) / 2
	res := make([]rune, rLen)
	for i := 0; i < rLen; i++ {
		n := r.Int31n(runes[i*2+1]-runes[i*2]+1) + runes[i*2]
		res[i] = n
	}
	return string(res[r.Intn(len(res))]), nil
}

func alternateGenerator(r *rand.Rand, parsedRegex *syntax.Regexp) (string, error) {
	sub := parsedRegex.Sub
	res := make([]string, 0, len(sub))
	for i := 0; i < len(sub); i++ {
		s, err := generate(r, sub[i])
		if err != nil {
			return "", err
		}
		res = append(res, s)
	}
	return res[r.Intn(len(sub))], nil
}

func captureGenerator(r *rand.Rand, parsedRegex *syntax.Regexp) (string, error) {
	return generate(r, parsedRegex.Sub[0])
}
