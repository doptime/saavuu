package data

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func (db *Ctx[k, v]) RPush(param ...v) (err error) {
	if val, err := db.toValueStrs(param); err != nil {
		return err
	} else {

		return db.Rds.RPush(db.Ctx, db.Key, val).Err()
	}
}
func (db *Ctx[k, v]) LPush(param ...v) (err error) {
	if val, err := db.toValueStrs(param); err != nil {
		return err
	} else {
		return db.Rds.LPush(db.Ctx, db.Key, val).Err()
	}
}
func (db *Ctx[k, v]) RPop() (ret v, err error) {
	cmd := db.Rds.RPop(db.Ctx, db.Key)
	if data, err := cmd.Bytes(); err != nil {
		return ret, err
	} else {
		return db.toValue(data)
	}
}
func (db *Ctx[k, v]) LPop() (ret v, err error) {
	cmd := db.Rds.LPop(db.Ctx, db.Key)
	if data, err := cmd.Bytes(); err != nil {
		return ret, err
	} else {
		return db.toValue(data)
	}
}
func (db *Ctx[k, v]) LRange(start, stop int64) (ret []v, err error) {
	cmd := db.Rds.LRange(db.Ctx, db.Key, start, stop)
	if data, err := cmd.Result(); err != nil {
		return ret, err
	} else {
		for _, val := range data {
			if v, err := db.toValue([]byte(val)); err != nil {
				return ret, err
			} else {
				ret = append(ret, v)
			}
		}
		return ret, nil
	}
}
func (db *Ctx[k, v]) LRem(count int64, param v) (err error) {
	if val, err := db.toValueStr(param); err != nil {
		return err
	} else {
		return db.Rds.LRem(db.Ctx, db.Key, count, val).Err()
	}
}
func (db *Ctx[k, v]) LSet(index int64, param v) (err error) {
	if val, err := db.toValueStr(param); err != nil {
		return err
	} else {
		return db.Rds.LSet(db.Ctx, db.Key, index, val).Err()
	}
}
func (db *Ctx[k, v]) BLPop(timeout time.Duration) (ret v, err error) {
	cmd := db.Rds.BLPop(db.Ctx, timeout, db.Key)
	if data, err := cmd.Result(); err != nil {
		return ret, err
	} else {
		return db.toValue([]byte(data[1]))
	}
}
func (db *Ctx[k, v]) BRPop(timeout time.Duration) (ret v, err error) {
	cmd := db.Rds.BRPop(db.Ctx, timeout, db.Key)
	if data, err := cmd.Result(); err != nil {
		return ret, err
	} else {
		return db.toValue([]byte(data[1]))
	}
}
func (db *Ctx[k, v]) BRPopLPush(destination string, timeout time.Duration) (ret v, err error) {
	cmd := db.Rds.BRPopLPush(db.Ctx, db.Key, destination, timeout)
	if data, err := cmd.Bytes(); err != nil {
		return ret, err
	} else {
		return db.toValue(data)
	}
}
func (db *Ctx[k, v]) LInsertBefore(pivot, param v) (err error) {
	var (
		pivotStr string
	)
	if val, err := db.toValueStr(param); err != nil {
		return err
	} else {
		if pivotStr, err = db.toValueStr(pivot); err != nil {
			return err
		} else {
			return db.Rds.LInsertBefore(db.Ctx, db.Key, pivotStr, val).Err()
		}
	}
}
func (db *Ctx[k, v]) LInsertAfter(pivot, param v) (err error) {
	var (
		pivotStr string
	)
	if val, err := db.toValueStr(param); err != nil {
		return err
	} else {
		if pivotStr, err = db.toValueStr(pivot); err != nil {
			return err
		} else {
			return db.Rds.LInsertAfter(db.Ctx, db.Key, pivotStr, val).Err()
		}
	}
}
func (db *Ctx[k, v]) Sort(sort *redis.Sort) (ret []v, err error) {
	cmd := db.Rds.Sort(db.Ctx, db.Key, sort)
	if data, err := cmd.Result(); err != nil {
		return ret, err
	} else {
		for _, val := range data {
			if v, err := db.toValue([]byte(val)); err != nil {
				return ret, err
			} else {
				ret = append(ret, v)
			}
		}
		return ret, nil
	}
}
func (db *Ctx[k, v]) LTrim(start, stop int64) (err error) {
	return db.Rds.LTrim(db.Ctx, db.Key, start, stop).Err()
}
func (db *Ctx[k, v]) LIndex(index int64) (ret v, err error) {
	cmd := db.Rds.LIndex(db.Ctx, db.Key, index)
	if data, err := cmd.Bytes(); err != nil {
		return ret, err
	} else {
		return db.toValue(data)
	}
}
func (db *Ctx[k, v]) LLen() (length int64, err error) {
	cmd := db.Rds.LLen(db.Ctx, db.Key)
	return cmd.Result()
}
