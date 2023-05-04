package permission

func init() {
	//if config.AppMode == config.AppModeFRAMEWROK
	go LoadPPermissionFromRedis()
}
