package main

import (
	"realWorld/global"
	"realWorld/initialize"
	"time"
)

func main() {
	//bindEnv方法里拥有阻塞以达到配置热加载
	//并发编程，防止热加载阻塞主进程
	time.Sleep(time.Second * 1)
	initialize.MustConfig()

	initialize.MustLoadGorm()
	initialize.AutoMigrate(global.DB)
	initialize.MustLoadZap()
	initialize.MustRunWindowServer()

}
