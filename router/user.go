package router

import (
	"github.com/gin-gonic/gin"
	"realWorld/api"
	"realWorld/middleware"
)

type UserGroup struct {
}

func (u *UserGroup) InitUserRouters(engine *gin.Engine) {

	//登录注册
	user := engine.Group("/api")
	{
		//用户身份验证
		user.POST("/users/login", api.UserLoginApi)
		//用户注册
		user.POST("/users", api.UserRegisterApi)
	}
	user_auth := engine.Group("/api")
	user_auth.Use(middleware.JwtMiddleware())
	{
		//获取当前用户
		user_auth.GET("/user", api.GetUserInfoApi)
		//更新用户
		user_auth.PUT("/user", api.PutUserInfoApi)
	}
}
