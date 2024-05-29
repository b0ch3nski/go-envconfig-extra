# go-envconfig-extra
[![license](https://img.shields.io/github/license/b0ch3nski/go-envconfig-extra)](LICENSE)
[![release](https://img.shields.io/github/v/release/b0ch3nski/go-envconfig-extra)](https://github.com/b0ch3nski/go-envconfig-extra/releases)
[![go.dev](https://pkg.go.dev/badge/github.com/b0ch3nski/go-envconfig-extra)](https://pkg.go.dev/github.com/b0ch3nski/go-envconfig-extra)
[![goreportcard](https://goreportcard.com/badge/github.com/b0ch3nski/go-envconfig-extra)](https://goreportcard.com/report/github.com/b0ch3nski/go-envconfig-extra)
[![issues](https://img.shields.io/github/issues/b0ch3nski/go-envconfig-extra)](https://github.com/b0ch3nski/go-envconfig-extra/issues)
[![sourcegraph](https://sourcegraph.com/github.com/b0ch3nski/go-envconfig-extra/-/badge.svg)](https://sourcegraph.com/github.com/b0ch3nski/go-envconfig-extra)

Extra tools extending usage of [go-envconfig](https://github.com/sethvargo/go-envconfig) library.

## install

```
go get github.com/b0ch3nski/go-envconfig-extra
```

## example

```go
import "github.com/b0ch3nski/go-envconfig-extra"

type Config struct {
	Password1 string `env:"PASS1,required" secret:"redact"`
	Password2 string `env:"PASS2,required" secret:"mask=4"`

	ArbitraryFile envconfigext.FileContent `env:"FILE"`
	Certificate   envconfigext.X509Cert    `env:"CERT"`
}

func (c Config) String() string {
	return envconfigext.StructFieldScan(c)
}
```
