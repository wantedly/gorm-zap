package gormzap

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/slashformotion/gorm-zap/testhelper"
)

var pool *testhelper.DockerPool

func TestMain(m *testing.M) {
	pool = testhelper.MustCreatePool()

	os.Exit(m.Run())
}

func Test_Logger_Postgres(t *testing.T) {
	fac, logs := observer.New(zap.DebugLevel)
	zapLogger := zap.New(fac)
	defer func() {
		err := zapLogger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	conn := pool.MustCreateDB(testhelper.DialectPostgres)
	defer conn.MustClose()

	now := time.Now()
	gorm.NowFunc = func() time.Time { return now }

	db, err := gorm.Open(conn.Dialect, conn.URL)
	if err != nil {
		panic(err)
	}

	type Post struct {
		Title, Body string
		CreatedAt   time.Time
	}
	db.AutoMigrate(&Post{})

	cases := []struct {
		run    func() error
		sql    string
		values []string
	}{
		{
			run: func() error { return db.Create(&Post{Title: "awesome"}).Error },
			sql: fmt.Sprintf(
				"INSERT INTO %q (%q,%q,%q) VALUES ($1,$2,$3) RETURNING %q.*",
				"posts", "title", "body", "created_at",
				"posts",
			),
			values: []string{"awesome", "", now.String()},
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
			sql: fmt.Sprintf(
				"SELECT * FROM %q  WHERE (%q.%q = $1) AND (%q.%q = $2) LIMIT 1",
				"posts", "posts", "title", "posts", "body",
			),
			values: []string{"awesome", "This is awesome post !"},
		},
	}

	db.SetLogger(New(zapLogger))
	db.LogMode(true)

	for _, c := range cases {
		err := c.run()
		if err != nil && err != gorm.ErrRecordNotFound {
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
