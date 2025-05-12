package router

import (
	"github.com/gin-gonic/gin"
	"realWorld/api"
	"realWorld/middleware"
)

type ProfilesGroup struct {
}

func (r *ProfilesGroup) InitProfileRouters(engine *gin.Engine) {
	profile := engine.Group("/api/profiles")
	profile.Use(middleware.JwtMiddleware())
	{
		profile.GET("/:username", api.GetUserProfileApi)
		profile.POST("/:username/follow", api.ProfileFollowApi)
		profile.DELETE("/:username/follow", api.ProfileUnFollowApi)
	}
}
