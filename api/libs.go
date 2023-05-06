package api

import (
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/data"
)

type InKeyLocker struct {
	Key        string
	DurationMs int64
}

var keyTimeLocker = data.New[bool]("TimeLocker")
var lastRemoveKeyLocker int64 = time.Now().UnixMicro()

var ApiKeyLocker = Api(func(req *InKeyLocker) (Locked bool, err error) {

	var (
		now     int64 = time.Now().UnixMicro()
		dueTime int64 = now + req.DurationMs
		score   float64
	)
	Locked = false
	if score, err = keyTimeLocker.ZScore(req.Key); err == nil && score > float64(now) {
		Locked = true
	}
	keyTimeLocker.ZAdd(redis.Z{Score: float64(dueTime), Member: req.Key})

	if now > lastRemoveKeyLocker {
		//remove key which expired every 1 hour
		lastRemoveKeyLocker += 36000 * 1000
		keyTimeLocker.ZRemRangeByScore("0", strconv.FormatInt(now, 10))
	}
	return Locked, nil
})
