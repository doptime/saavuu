package key

import (
	"fmt"
	"time"
)

func YMD(Key string, tm time.Time) string {
	//year is 4 digits, month is 2 digits, day is 2 digits
	return fmt.Sprintf("%sYMD:%04v%02v%02v", Key, tm.Year(), int(tm.Month()), tm.Day())
}
func YM(Key string, tm time.Time) string {
	//year is 4 digits, month is 2 digits
	return fmt.Sprintf("%sYM:%04v%02v", Key, tm.Year(), int(tm.Month()))
}
func Y(Key string, tm time.Time) string {
	//year is 4 digits
	return fmt.Sprintf("%sY:%04v", Key, tm.Year())
}
func YW(Key string, tm time.Time) string {
	tm = tm.UTC()
	isoYear, isoWeek := tm.ISOWeek()
	//year is 4 digits, week is 2 digits
	return fmt.Sprintf("%sYW:%04v%02v", Key, isoYear, isoWeek)
}
func Multiple(key string, fields ...string) string {
	//	concacate all fields with ':'
	strAll := key
	for _, field := range fields {
		strAll += ":" + field
	}
	return strAll
}
