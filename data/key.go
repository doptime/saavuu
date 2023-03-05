package data

import (
	"fmt"
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
func (db *Ctx) ConcatedKey(fields ...interface{}) string {
	//	concacate all fields with ':'
	strAll := string(db.Key)
	for _, field := range fields {
		//convert field to string,field may be int, float, string, etc.
		strAll += fmt.Sprintf(":%v", field)
	}
	return strAll
}

func (db *Ctx) Concat(fields ...interface{}) *Ctx {
	//	concacate all fields with ':'
	strAll := string(db.Key)
	for _, field := range fields {
		//convert field to string,field may be int, float, string, etc.
		strAll += fmt.Sprintf(":%v", field)
	}
	return &Ctx{db.Ctx, db.Rds, db.ConcatedKey()}
}
