package api

import (
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/data"
)

type InTimeLocker struct {
	Key            string
	LockDurationMs int64
}
type OutTimeLocker struct {
	Locked bool
}

var keyTimeLocker = data.New[bool]("TimeLocker")
var lastRemoveTimeLockerExpiredTime int64 = time.Now().UnixNano() / 1000000 //ms

var ApiTimeLocker = Api(func(req *InTimeLocker) (out *OutTimeLocker, err error) {

	var (
		now     int64 = time.Now().UnixNano() / 1000000
		dueTime int64 = now + req.LockDurationMs
		score   float64
	)
	out = &OutTimeLocker{Locked: false}
	if score, err = keyTimeLocker.ZScore(req.Key); err == nil && score > float64(now) {
		out.Locked = true
	}
	keyTimeLocker.ZAdd(redis.Z{Score: float64(dueTime), Member: req.Key})

	if now > lastRemoveTimeLockerExpiredTime {
		//remove key which expired every 1 hour
		lastRemoveTimeLockerExpiredTime += 36000 * 1000
		keyTimeLocker.ZRemRangeByScore("0", strconv.FormatInt(now, 10))
	}
	return out, nil
})
