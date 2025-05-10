package api

import (
	"github.com/gin-gonic/gin"
	"realWorld/model"
	"realWorld/model/resp"
	"realWorld/service"
	"realWorld/utils"
)

func GetUserProfile(c *gin.Context) {
	//获取搜索用户名字
	target := c.Param("username")

	//获取登录用户的名字
	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage("用户信息出错，请重新登录", c)
		return
	}
	//业务处理
	targetId := getTargetId(target)
	if targetId == 0 {
		resp.FailWithMessage("没有该用户", c)
		return
	}

	data := &model.Follower{
		UserId:     claims.Id,
		FollowerId: targetId,
	}
	userInfo := service.ProfileServiceApp
	info, err := userInfo.GetTagetUserInfo(data)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	resp.OkWithData(info, c)

}

func ProfileFollow(c *gin.Context) {
	//获取搜索用户名字
	target := c.Param("username")

	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage("用户信息出错，请重新登录", c)
		return
	}
	targetId := getTargetId(target)
	if targetId == 0 {
		resp.FailWithMessage("没有该用户", c)
		return
	}

	data := &model.Follower{
		UserId:     claims.Id,
		FollowerId: targetId,
	}
	userInfo := service.ProfileServiceApp
	user, err := userInfo.FollowUser(data)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	resp.OkWithData(user, c)
}

func ProfileUnFollow(c *gin.Context) {
	//获取搜索用户名字
	target := c.Param("username")

	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage("用户信息出错，请重新登录", c)
		return
	}
	userInfo := service.ProfileServiceApp
	targetId := getTargetId(target)
	if targetId == 0 {
		resp.FailWithMessage("没有该用户", c)
		return
	}
	data := &model.Follower{
		UserId:     claims.Id,
		FollowerId: targetId,
	}
	user, err := userInfo.UnFollowUser(data)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	resp.OkWithData(user, c)
}

// 私有方法 获取搜索用户id根据用户名
func getTargetId(username string) uint {
	user := service.UserService{}
	info, err := user.GetUserInfo(username)
	if err != nil {
		return 0
	}
	return info.User.Id
}
