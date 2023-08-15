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