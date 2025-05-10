package initialize

import (
	"flag"
	"fmt"
	"os"
	"realWorld/global"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 创建`MustConfig`函数来读取配置文件
func MustConfig() {
	// 根据命令行参数、环境变量、Gin Mode实现配置文件切换
	var config string
	flag.StringVar(&config, "c", "", "指定配置文件路径")
	flag.Parse()
	if config == "" { // 判断命令行参数是否指定了配置文件
		if configEnv := os.Getenv("CONFIG"); configEnv == "" { // 读取环境变量
			switch gin.Mode() {
			case gin.DebugMode:
				config = "config.yaml"
			case gin.TestMode:
				config = "config.test.yaml"
			case gin.ReleaseMode:
				config = "config.release.yaml"
			}
			fmt.Printf("您正在使用gin模式的%s环境名称, 配置文件的路径为%s\n", gin.Mode(), config)
		} else {
			config = configEnv
			fmt.Printf("您正在使用环境变量, 配置文件的路径为%s\n", configEnv)
		}
	} else {
		fmt.Printf("您正在使用命令行的-c参数传递的值, 配置文件的路径为%s\n", config)
	}

	// 绑定到结构体中
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = v.Unmarshal(&global.CONFIG)
	if err != nil {
		panic(err)
	}
	v.WatchConfig()

	// 配置热加载功能
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件发生变更: ", in.Name)
		err = v.Unmarshal(&global.CONFIG)
		if err != nil {
			panic(err)
		}
		fmt.Println(global.CONFIG)
	})
}
