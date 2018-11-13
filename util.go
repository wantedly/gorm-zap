package gormzap

import (
	"fmt"
	"regexp"
	"time"
	"unicode"
)

var (
	sqlRegexp                = regexp.MustCompile(`\?`)
	numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)
)

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func getFormattedValues(values []interface{}) []string {
	rawValues := values[4].([]interface{})
	formattedValues := make([]string, 0, len(rawValues))
	for _, value := range rawValues {
		switch v := value.(type) {
		case time.Time:
			formattedValues = append(formattedValues, fmt.Sprint(v))
		case []byte:
			if str := string(v); isPrintable(str) {
				formattedValues = append(formattedValues, fmt.Sprint(str))
			} else {
				formattedValues = append(formattedValues, "<binary>")
			}
		default:
			str := "NULL"
			if v != nil {
				str = fmt.Sprint(v)
			}
			formattedValues = append(formattedValues, str)
		}
	}
	return formattedValues
}

func getFormatSQL(values []interface{}, ret []string) string {
	var sql string

	// differentiate between $n placeholders or else treat like ?
	if numericPlaceHolderRegexp.MatchString(values[3].(string)) {
		sql = values[3].(string)
		for index, value := range ret {
			placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
			sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
		}
	} else {
		formattedValuesLength := len(ret)
		for index, value := range sqlRegexp.Split(values[3].(string), -1) {
			sql += value
			if index < formattedValuesLength {
				sql += ret[index]
			}
		}
	}

	return sql
}

func getSource(values []interface{}) string {
	return fmt.Sprint(values[1])
}

func getDuration(values []interface{}) time.Duration {
	return values[2].(time.Duration)
}
