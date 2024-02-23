package permission

import (
	"fmt"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/rs/zerolog/log"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/data"
)

var rdsPermit = data.New[string, string]("_permissions")
var permitmap cmap.ConcurrentMap[string, bool] = cmap.New[bool]()

// this version of IsPermitted is design for fast searching & modifying
func IsPermitted(dataKey string, operation string) (ok bool) {
	var (
		autoPermit                bool   = config.Cfg.Data.AutoAuth
		keyAllowed, keyDisAllowed string = fmt.Sprintf("%s::%s::on", dataKey, operation), fmt.Sprintf("%s::%s::off", dataKey, operation)
	)
	if _, ok := permitmap.Get(keyAllowed); ok {
		return true
	}
	if _, ok := permitmap.Get(keyDisAllowed); ok {
		return false
	}
	if autoPermit {
		permitmap.Set(keyAllowed, true)
		rdsPermit.HSet(keyAllowed, time.Now().Format("2006-01-02 15:04:05"))
	}
	return autoPermit
}

var ConfigurationLoaded bool = false

func LoadPermissionTable() {
	var (
		keys []string
		err  error
	)

	if keys, err = rdsPermit.Keys(); !ConfigurationLoaded {
		if err != nil {
			log.Warn().AnErr("Step2.1: start permission loading from redis failed", err).Send()
		} else {
			log.Info().Msg("Step2.2: start permission loaded from redis")
		}
		ConfigurationLoaded = true
	}
	for _, key := range keys {
		if _, ok := permitmap.Get(key); !ok {
			permitmap.Set(key, true)
		}
	}
	go func() {
		time.Sleep(time.Minute)
		LoadPermissionTable()
	}()
}

func init() {
	LoadPermissionTable()
}
