package permission

import "github.com/yangkequn/saavuu/config"

func init() {
	if config.AppMode == config.AppModeFRAMEWROK {
		go LoadPPermissionFromRedis()
		//init permission
	}
}
