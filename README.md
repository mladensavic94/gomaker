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
BenchmarkRandFill-20                              136150              8866 ns/op            6143 B/op         15 allocs/op
BenchmarkRandFill_WithPreloadMapping-20           137984              8210 ns/op            5760 B/op          9 allocs/op
BenchmarkRegexFill-20                             113206             10401 ns/op            7518 B/op         58 allocs/op
BenchmarkRegexFill_WithPreloadMapping-20          122317              9481 ns/op            7066 B/op         47 allocs/op
BenchmarkFuncFill-20                              148395              8459 ns/op            5812 B/op         11 allocs/op
BenchmarkFuncFill_WithPreloadMapping-20           150086              7505 ns/op            5428 B/op          5 allocs/op
PASS
ok      gomaker 7.665s
```

## TODO
1. use map instead of reflect search every time -
2. return value instead of fill
3. generics
4. rel options