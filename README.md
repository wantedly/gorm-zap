# gormzap
[![Build Status](https://travis-ci.org/wantedly/gorm-zap.svg?branch=master)](https://travis-ci.org/wantedly/gorm-zap)
[![codecov](https://codecov.io/gh/wantedly/gorm-zap/branch/master/graph/badge.svg)](https://codecov.io/gh/wantedly/gorm-zap)
[![GoDoc](https://godoc.org/github.com/wantedly/gorm-zap?status.svg)](https://godoc.org/github.com/wantedly/gorm-zap)
[![license](https://img.shields.io/github/license/wantedly/gorm-zap.svg)](./LICENSE)

Alternative logging with [zap](https://github.com/uber-go/zap) for [GORM](http://jinzhu.me/gorm) ⚡️

In comparison to gorm's default logger, `gormzap` is faster, reflection free, low allocations and no regex compilations.


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


## Performance
According to our benchmark, `gormzap` makes DB operations at least 5% faster and reduce object allocations.

### Simple insert query

| Logger | Time | Object Allocated |
| :--- | :---: | :---: |
| default | 187940 ns/op | 494 allocs/op |
| gormzap | 185383 ns/op | 475 allocs/op |

### Simple select query

| Logger | Time | Object Allocated |
| :--- | :---: | :---: |
| default | 169361 ns/op | 531 allocs/op |
| gormzap | 151304 ns/op | 519 allocs/op |

### Simple select query with 10 placeholders

| Logger | Time | Object Allocated |
| :--- | :---: | :---: |
| default | 200632 ns/op | 720 allocs/op |
| gormzap | 190732 ns/op | 645 allocs/op |

### Simple select query with 100 placeholders

| Logger | Time | Object Allocated |
| :--- | :---: | :---: |
| default | 444513 ns/op | 1723 allocs/op |
| gormzap | 263098 ns/op | 1101 allocs/op |
