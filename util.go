package gormzap

import (
	"fmt"
	"time"
	"unicode"
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
	formattedValues := []string{}
	for _, value := range values[4].([]interface{}) {
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

func getSource(values []interface{}) string {
	return fmt.Sprint(values[1])
}

func getDuration(values []interface{}) time.Duration {
	return values[2].(time.Duration)
}
