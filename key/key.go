package key

import (
	"fmt"
	"time"
)

type RedisKey string

func (Key RedisKey) YMD(tm time.Time) string {
	//year is 4 digits, month is 2 digits, day is 2 digits
	return fmt.Sprintf("%sYMD:%04v%02v%02v", Key, tm.Year(), int(tm.Month()), tm.Day())
}
func (Key RedisKey) YM(tm time.Time) string {
	//year is 4 digits, month is 2 digits
	return fmt.Sprintf("%sYM:%04v%02v", Key, tm.Year(), int(tm.Month()))
}
func (Key RedisKey) Y(tm time.Time) string {
	//year is 4 digits
	return fmt.Sprintf("%sY:%04v", Key, tm.Year())
}
func (Key RedisKey) YW(tm time.Time) string {
	tm = tm.UTC()
	isoYear, isoWeek := tm.ISOWeek()
	//year is 4 digits, week is 2 digits
	return fmt.Sprintf("%sYW:%04v%02v", Key, isoYear, isoWeek)
}
func (Key RedisKey) Concat(fields ...interface{}) string {
	//	concacate all fields with ':'
	strAll := string(Key)
	for _, field := range fields {
		//convert field to string,field may be int, float, string, etc.
		strAll += fmt.Sprintf(":%v", field)
	}
	return strAll
}
