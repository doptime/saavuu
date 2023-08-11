package data

import (
	"github.com/redis/go-redis/v9"
)

// sacn key by pattern
func (db *Ctx[k, v]) Scan(match string, cursor uint64, count int64) (keys []string, err error) {
	var (
		cmd   *redis.ScanCmd
		_keys []string
	)
	//scan all keys
	for {

		if cmd = db.Rds.Scan(db.Ctx, cursor, match, count); cmd.Err() != nil {
			return nil, cmd.Err()
		}
		if _keys, cursor, err = cmd.Result(); err != nil {
			return nil, err
		}
		keys = append(keys, _keys...)
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}
