package permission

import "github.com/yangkequn/saavuu/config"

func init() {
	if config.AppMode == config.AppModeFRAMEWROK {
		go LoadGetPermissionFromRedis()
		go LoadPutPermissionFromRedis()
		go LoadDelPermissionFromRedis()
		//init permission
	}
}
