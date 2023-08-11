package data

func (db *Ctx[k, v]) RPush(param v) (err error) {
	if val, err := db.toValueStr(param); err != nil {
		return err
	} else {

		return db.Rds.RPush(db.Ctx, db.Key, val).Err()
	}
}
func (db *Ctx[k, v]) LSet(index int64, param v) (err error) {
	if val, err := db.toValueStr(param); err != nil {
		return err
	} else {
		return db.Rds.LSet(db.Ctx, db.Key, index, val).Err()
	}
}
func (db *Ctx[k, v]) LGet(index int64) (ret v, err error) {
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
