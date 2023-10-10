# gomaker
![Build](https://github.com/mladensavic94/gomaker/actions/workflows/go.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/mladensavic94/gomaker/badge.svg?branch=main)](https://coveralls.io/github/mladensavic94/gomaker?branch=main)

Fill structs with random data.

```go
type dummy struct {
    Id          int64  `gomaker:"rand[1;10;1]"`
    DummyString string `gomaker:"rand"`
}

maker := gomaker.New()
d := dummy{}
err := maker.Fill(&d)
```

## Benchmarks
```shell
go test -bench=^Bench -count 1 -run=^# -benchmem
goos: windows
goarch: amd64
pkg: gomaker
cpu: 12th Gen Intel(R) Core(TM) i9-12900H
BenchmarkRandFill-20              138298              8621 ns/op            6144 B/op         15 allocs/op
BenchmarkRegexFill-20             123716              9518 ns/op            7443 B/op         52 allocs/op
BenchmarkFuncFill-20              148338              7958 ns/op            5812 B/op         11 allocs/op
PASS
ok      gomaker 3.862s
```

## TODO
1. use map instead of reflect search every time -
2. return value instead of fill
3. generics
4. rel options