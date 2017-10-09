package gormzap

import (
	"database/sql/driver"
	"io/ioutil"
	stdlog "log"
	"testing"

	"github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func Benchmark_WithTestDB(b *testing.B) {
	// https://github.com/uber-go/zap/blob/35aad584952c3e7020db7b839f6b102de6271f89/benchmarks/zap_test.go#L106-L116
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeDuration = zapcore.NanosDurationEncoder
	ec.EncodeTime = zapcore.EpochNanosTimeEncoder
	enc := zapcore.NewJSONEncoder(ec)
	zapLogger := zap.New(zapcore.NewCore(
		enc,
		&zaptest.Discarder{},
		zap.DebugLevel,
	))

	defer zapLogger.Sync()

	type Post struct {
		ID    int
		Title string
		Body  string
	}

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (driver.Result, error) {
		return testdb.NewResult(1, nil, 1, nil), nil
	})
	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		return testdb.RowsFromCSVString([]string{"title", "body"}, `"awesome","This is an awesome post"`), nil
	})

	setupDB := func() *gorm.DB {
		db, err := gorm.Open("testdb", "")
		if err != nil {
			b.Fatal(err)
		}
		db.AutoMigrate(&Post{})
		db.LogMode(true)
		return db
	}

	benchInsert := func(b *testing.B, db *gorm.DB) {
		post := &Post{Title: "awesome", Body: "This is an awesome post"}
		b.ResetTimer()
		for i := 1; i <= b.N; i++ {
			db.Create(post)
		}
	}
	benchSelectByID := func(b *testing.B, db *gorm.DB) {
		b.ResetTimer()
		for i := 1; i <= b.N; i++ {
			db.Model(&Post{}).Where(&Post{ID: i}).Find(&[]*Post{})
		}
	}
	benchSelectByIDs := func(b *testing.B, db *gorm.DB, n int) {
		ids := make([]int, n)
		for i := 1; i <= n; i++ {
			ids = append(ids, i)
		}
		b.ResetTimer()
		for i := 1; i <= b.N; i++ {
			db.Model(&Post{}).Where("id in (?)", ids).Find(&[]*Post{})
		}
	}

	b.Run("default", func(b *testing.B) {
		db := setupDB()
		defer db.Close()
		// https://github.com/jinzhu/gorm/blob/3a9e91ab372120a0e35b518430255308e3d8d5ea/logger.go#L16
		db.SetLogger(gorm.Logger{LogWriter: stdlog.New(ioutil.Discard, "\r\n", 0)})

		b.ResetTimer()
		b.Run("insert post", func(b *testing.B) { benchInsert(b, db) })
		b.Run("select by ID", func(b *testing.B) { benchSelectByID(b, db) })
		b.Run("select by 10 IDs", func(b *testing.B) { benchSelectByIDs(b, db, 10) })
		b.Run("select by 100 IDs", func(b *testing.B) { benchSelectByIDs(b, db, 100) })
	})

	b.Run("gormzap", func(b *testing.B) {
		db := setupDB()
		defer db.Close()
		db.SetLogger(New(zapLogger))

		b.ResetTimer()
		b.Run("insert post", func(b *testing.B) { benchInsert(b, db) })
		b.Run("select by ID", func(b *testing.B) { benchSelectByID(b, db) })
		b.Run("select by 10 IDs", func(b *testing.B) { benchSelectByIDs(b, db, 10) })
		b.Run("select by 100 IDs", func(b *testing.B) { benchSelectByIDs(b, db, 100) })
	})
}
