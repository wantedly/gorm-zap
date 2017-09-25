package zapgorm

import (
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func Test_Logger(t *testing.T) {
	fac, logs := observer.New(zap.DebugLevel)
	zapLogger := zap.New(fac)
	defer func() {
		err := zapLogger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	db, err := gorm.Open("testdb", "")
	if err != nil {
		panic(fmt.Errorf("unexpected error: %v", err))
	}

	db.SetLogger(FromZap(zapLogger))
	db.LogMode(true)

	type Post struct{ Title, Body string }

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (driver.Result, error) {
		return testdb.NewResult(1, nil, 1, nil), nil
	})
	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		return testdb.RowsFromCSVString([]string{"title"}, "awesome"), nil
	})

	cases := []struct {
		run    func() error
		sql    string
		values []string
	}{
		{
			run:    func() error { return db.Create(&Post{Title: "awesome"}).Error },
			sql:    "INSERT INTO \"posts\" (\"title\",\"body\") VALUES (?,?)",
			values: []string{"awesome", ""},
		},
		{
			run:    func() error { return db.Model(&Post{}).Find(&[]*Post{}).Error },
			sql:    "SELECT * FROM \"posts\"  ",
			values: []string{},
		},
		{
			run: func() error {
				return db.Where(&Post{Title: "awesome", Body: "This is awesome post !"}).First(&Post{}).Error
			},
			sql:    "SELECT * FROM \"posts\"  WHERE (\"title\" = ?) AND (\"body\" = ?) LIMIT 1",
			values: []string{"awesome", "This is awesome post !"},
		},
	}

	for _, c := range cases {
		err := c.run()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		entries := logs.TakeAll()

		if got, want := len(entries), 1; got != want {
			t.Errorf("Logger logged %d items, want %d items", got, want)
		}

		fieldByName := entries[0].ContextMap()

		if got, want := fieldByName["sql"].(string), c.sql; got != want {
			t.Errorf("Logged sql was %q, want %q", got, want)
		}

		if got, want := len(fieldByName["values"].([]interface{})), len(c.values); got != want {
			t.Errorf("Logged values has %d items, want %d items", got, want)
		}

		for i, want := range c.values {
			got := fieldByName["values"].([]interface{})[i].(string)
			if got != want {
				t.Errorf("Logged values at %d was %v, want %v", i, got, want)
			}
		}
	}
}
