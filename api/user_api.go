package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"realWorld/global"
	"realWorld/middleware"
	"realWorld/model"
	"realWorld/model/req"
	"realWorld/model/resp"
	"realWorld/service"
	"realWorld/utils"
)

func UserLoginApi(c *gin.Context) {
	var userReq req.UserAuth
	//获取请求体值
	err := c.ShouldBindJSON(&userReq)
	fmt.Println(userReq)
	if err != nil {
		resp.FailWithMessage(utils.Translate(err), c)
		return
	}
	//排除postman的接口错误
	if userReq.User.Email == "{{EMAIL}}" || userReq.User.Password == "{{PASSWORD}}" {
		resp.FailWithMessage("传值为空", c)
		return
	}
	//合法性认证
	//处理业务
	//todo
	//userReq.User.Password = middleware.Md5Decode(userReq.User.Password)
	//创建一个实体类
	userModel := service.UserServiceApp
	login, err := userModel.Login(userReq.User)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)

		//写入日志
		global.Logger.Warn("登录失败"+err.Error(), zap.String("service", "login"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	baseclaim := model.BaseClaims{
		Id:       login.User.Id,
		Username: login.User.Username,
		Email:    login.User.Email,
	}

	genJwt := utils.NewJwt()
	tokenstr := genJwt.CreateClaims(baseclaim)
	//生成token
	token, err := genJwt.GenerateToken(&tokenstr)
	if err != nil {
		resp.FailWithMessage("token生成错误", c)

		//写入日志
		global.Logger.Warn("token生成失败"+err.Error(), zap.String("service", "gentoken"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	fmt.Println(token)
	login.User.Token = token
	//成功返回结果
	resp.OkWithData(login, c)

}

func UserRegisterApi(c *gin.Context) {
	var userReq req.UserRegister
	err := c.ShouldBindJSON(&userReq)
	if err != nil {
		resp.FailWithMessage(utils.Translate(err), c)
		return
	}
	if userReq.User.Email == "{{EMAIL}}" || userReq.User.Password == "{{PASSWORD}}" || userReq.User.Username == "{{USERNAME}}" {
		resp.FailWithMessage("传值为空", c)
		return
	}
	//合法性认证
	//密码加入随机种子并加密 todo
	userReq.User.Password = middleware.Md5Decode(userReq.User.Password)

	//处理业务
	userModel := service.UserServiceApp
	register, err := userModel.Register(userReq.User)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("注册失败"+err.Error(), zap.String("service", "register"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	baseclaim := model.BaseClaims{
		Username: register.User.Username,
		Email:    register.User.Email,
	}
	fmt.Println(baseclaim)
	genJwt := utils.NewJwt()
	tokenstr := genJwt.CreateClaims(baseclaim)
	//生成token
	token, err := genJwt.GenerateToken(&tokenstr)
	if err != nil {
		resp.FailWithMessage("token生成错误", c)
		//写入日志
		global.Logger.Warn("token生成失败"+err.Error(), zap.String("service", "gentoken"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	register.User.Token = token
	//成功返回结果
	resp.OkWithData(register, c)
}

// 获取当前用户
func GetUserInfoApi(c *gin.Context) {
	//	从token中获取信息
	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	//	先从redis中查找我的信息

	//	从mysql中查找
	userModel := service.UserServiceApp
	info, err := userModel.GetUserInfo(claims.Username)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("获取用户信息"+err.Error(), zap.String("service", "getuserinfo"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}

	resp.OkWithData(info, c)
}

// 更改用户的个人信息
// 更改用户的个人信息
func PutUserInfoApi(c *gin.Context) {
	// 根据token获取用户信息
	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}

	// 绑定请求数据
	var user req.ModifyUser
	if err = c.ShouldBindJSON(&user); err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}

	// 验证请求数据
	validate := validator.New()
	if err = validate.Struct(user); err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}

	// 调用服务层更新用户信息
	userModel := service.UserServiceApp
	updatedUser, err := userModel.ModifyUserInfo(&user.User, claims.Username)
	if err != nil {
		resp.FailWithMessage(utils.Translate(err), c)
		return
	}

	// 检查是否需要重新生成token
	if needsNewToken(updatedUser, claims) {
		token, err := generateNewToken(*updatedUser)
		if err != nil {
			global.Logger.Warn("更新数据失败", zap.Error(err), zap.String("service", "putuserinfo"), zap.Int("port", global.CONFIG.Server.Port))
			resp.FailWithMessage("token生成错误", c)
			return
		}
		resp.OkWithData(token, c)
		return
	}

	// 更新成功，无需重新生成token
	resp.OkWithMessage("更新成功", c)
}

// 检查是否需要重新生成token
func needsNewToken(updatedUser *resp.UserResp, claims *model.GoShopClaims) bool {
	return updatedUser.User.Email != claims.Email || updatedUser.User.Username != claims.Username
}

// 生成新的token
func generateNewToken(user resp.UserResp) (string, error) {
	baseClaim := model.BaseClaims{
		Id:       user.User.Id,
		Username: user.User.Username,
		Email:    user.User.Email,
	}
	genJwt := utils.NewJwt()
	tokenStr := genJwt.CreateClaims(baseClaim)
	return genJwt.GenerateToken(&tokenStr)
}
