# gomaker
![Build](https://github.com/mladensavic94/gomaker/actions/workflows/go.yml/badge.svg)

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