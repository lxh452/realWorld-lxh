package initialize

import (
	"fmt"
	"realWorld/global"
	"realWorld/router"

	"github.com/gin-gonic/gin"
)

// 创建`MustRunWindowServer`来实现`gin`服务的运行
func MustRunWindowServer() {
	engine := gin.Default()

	userGroup := router.UserGroup{}
	userGroup.InitUserRouters(engine)

	ArticleGroup := router.ArticleGroup{}
	ArticleGroup.InitArticleRouters(engine)

	profilesGroup := router.ProfilesGroup{}
	profilesGroup.InitProfileRouters(engine)

	address := fmt.Sprintf(":%d", global.CONFIG.Server.Port)
	fmt.Println("启动服务器，监听端口：", address)
	if err := engine.Run(address); err != nil {
		panic(err)
	}

}
