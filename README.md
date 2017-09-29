# gormzap
[![Build Status](https://travis-ci.org/wantedly/gorm-zap.svg?branch=master)](https://travis-ci.org/wantedly/gorm-zap)
[![codecov](https://codecov.io/gh/wantedly/gorm-zap/branch/master/graph/badge.svg)](https://codecov.io/gh/wantedly/gorm-zap)
[![GoDoc](https://godoc.org/github.com/wantedly/gorm-zap?status.svg)](https://godoc.org/github.com/wantedly/gorm-zap)
[![license](https://img.shields.io/github/license/wantedly/gorm-zap.svg)](./LICENSE)

Alternative logging with zap for GORM ⚡️

## Example

```go
package main

import (
	"github.com/jinzhu/gorm"
	"github.com/wantedly/gorm-zap"
)

const (
	databaseURL = "postgres://postgres:@localhost/gormzap?sslmode=disable"
)

func main() {
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open("postgres", databaseURL)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	db.SetLogger(gormzap.New(logger))

	// ...
}
```
