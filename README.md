# gormzap
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
