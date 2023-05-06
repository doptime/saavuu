package api

import (
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/data"
)

type InLockKey struct {
	Key        string
	DurationMs int64
}

var keyTimeLocker = data.New[bool]("TimeLocker")
var removeCounter int64 = 90

var ApiLockKey = Api(func(req *InLockKey) (ok bool, err error) {

	var (
		now     int64 = time.Now().UnixMicro()
		dueTime int64 = now + req.DurationMs
		score   float64
	)
	ok = true
	if score, err = keyTimeLocker.ZScore(req.Key); err == nil && score > float64(now) {
		ok = false
	}
	keyTimeLocker.ZAdd(redis.Z{Score: float64(dueTime), Member: req.Key})

	//auto remove expired keys
	removeCounter += 1
	if removeCounter > 100 {
		removeCounter = 0
		go keyTimeLocker.ZRemRangeByScore("0", strconv.FormatInt(now, 10))
	}
	return ok, nil
})
