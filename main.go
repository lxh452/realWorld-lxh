package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"realWorld/global"
	"realWorld/initialize"
	"time"
)

func main() {
	//bindEnv方法里拥有阻塞以达到配置热加载
	//并发编程，防止热加载阻塞主进程
	time.Sleep(time.Second * 1)
	initialize.MustConfig()
	go func() { // 启动 pprof HTTP 服务
		pprofAddress := ":6060" // pprof 监听端口
		fmt.Println("启动 pprof 服务，监听端口：", pprofAddress)
		if err := http.ListenAndServe(pprofAddress, nil); err != nil {
			fmt.Println("pprof 服务启动失败:", err)
		}
	}()
	initialize.MustLoadGorm()
	initialize.AutoMigrate(global.DB)
	initialize.MustLoadZap()
	initialize.MustRunWindowServer()

}
