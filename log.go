package gormzap

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type log struct {
	occurredAt time.Time
	source     string
	duration   int64
	sql        string
	values     []string
	other      []string
}

func (l *log) toZapFields() []zapcore.Field {
	return []zapcore.Field{
		zap.String("occurredAt", fmt.Sprint(l.occurredAt)),
		zap.String("source", l.source),
		zap.Int64("duration", l.duration),
		zap.String("sql", l.sql),
		zap.Strings("values", l.values),
		zap.Strings("other", l.other),
	}
}

func createLog(values []interface{}) *log {
	ret := &log{}
	ret.occurredAt = gorm.NowFunc()

	if len(values) > 1 {
		var level = values[0]
		ret.source = getSource(values)

		if level == "sql" {
			ret.duration = getDuration(values)
			ret.values = getFormattedValues(values)
			ret.sql = values[3].(string)
		} else {
			ret.other = append(ret.other, fmt.Sprint(values[2:]))
		}
	}

	return ret
}
