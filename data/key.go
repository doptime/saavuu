package data

import (
	"fmt"
	"strings"
	"time"
)

func (db *Ctx) YMD(tm time.Time) *Ctx {
	//year is 4 digits, month is 2 digits, day is 2 digits
	return &Ctx{db.Ctx, db.Rds, fmt.Sprintf("%sYMD:%04v%02v%02v", db.Key, tm.Year(), int(tm.Month()), tm.Day())}
}
func (db *Ctx) YM(tm time.Time) *Ctx {
	//year is 4 digits, month is 2 digits
	return &Ctx{db.Ctx, db.Rds, fmt.Sprintf("%sYM:%04v%02v", db.Key, tm.Year(), int(tm.Month()))}
}
func (db *Ctx) Y(tm time.Time) *Ctx {
	//year is 4 digits
	return &Ctx{db.Ctx, db.Rds, fmt.Sprintf("%sY:%04v", db.Key, tm.Year())}
}
func (db *Ctx) YW(tm time.Time) *Ctx {
	tm = tm.UTC()
	isoYear, isoWeek := tm.ISOWeek()
	//year is 4 digits, week is 2 digits
	return &Ctx{db.Ctx, db.Rds, fmt.Sprintf("%sYW:%04v%02v", db.Key, isoYear, isoWeek)}
}
func ConcatedKeys(fields ...interface{}) string {
	//	concacate all fields with ':'
	var key string
	for _, field := range fields {
		key += fmt.Sprintf("%v:", field)
	}
	//	remove the last ':'
	if len(key) == 0 {
		return ""
	}
	return key[:len(key)-1]
}

func (db *Ctx) Concat(fields ...interface{}) *Ctx {
	//for each field ,it it's type if float64 or float32,but it's value is integer,then convert it to int
	for i, field := range fields {
		if f64, ok := field.(float64); ok && f64 == float64(int64(f64)) {
			fields[i] = int64(field.(float64))
		} else if f32, ok := field.(float32); ok && f32 == float32(int32(f32)) {
			fields[i] = int32(field.(float32))
		}
	}
	//implete logic of  return &Ctx{db.Ctx, db.Rds, fmt.Sprintf("%s:%v", db.Key, ConcatedKeys(fields...))}
	//but ,do not use recursion
	results := make([]string, 0, len(fields)+1)
	results = append(results, db.Key)
	for _, field := range fields {
		results = append(results, fmt.Sprintf("%v", field))
	}
	return &Ctx{db.Ctx, db.Rds, strings.Join(results, ":")}
}
