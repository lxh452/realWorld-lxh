package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
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
		log.Println("绑定请求体失败", err.Error())
		resp.FailWithMessage(err.Error(), c)
		return
	}
	//排除postman的接口错误
	if userReq.User.Email == "{{EMAIL}}" || userReq.User.Password == "{{PASSWORD}}" {
		resp.FailWithMessage("传值为空", c)
		return
	}
	//合法性认证
	//处理业务
	//创建一个实体类
	userModel := service.UserServiceApp
	login, err := userModel.Login(userReq.User)
	if err != nil {
		log.Println(err.Error())
		resp.FailWithMessage(err.Error(), c)
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
		log.Println("绑定请求体失败", err.Error())
		resp.FailWithMessage(err.Error(), c)
		return
	}
	if userReq.User.Email == "{{EMAIL}}" || userReq.User.Passwd == "{{PASSWORD}}" || userReq.User.Username == "{{USERNAME}}" {
		resp.FailWithMessage("传值为空", c)
		return
	}
	//合法性认证
	//处理业务
	userModel := service.UserServiceApp
	register, err := userModel.Register(userReq.User)
	if err != nil {
		log.Println(err.Error())
		resp.FailWithMessage(err.Error(), c)
		return
	}
	baseclaim := model.BaseClaims{
		Id:       register.User.Id,
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
		return
	}
	register.User.Token = token
	//成功返回结果
	resp.OkWithData(register, c)
}

// 获取当前用户
func GetUserInfo(c *gin.Context) {
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
		return
	}

	resp.OkWithData(info, c)
}

// 更改用户的个人信息
func PutUserInfo(c *gin.Context) {
	//根据token获取用户信息
	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	var user req.ModifyUser
	if err = c.ShouldBindJSON(&user); err != nil {
		resp.FailWithMessage(err.Error(), c)
		fmt.Println(user)
		return
	}

	validate := validator.New()
	if err = validate.Struct(user); err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	userModel := service.UserServiceApp
	_, err = userModel.ModifyUserInfo(&user.User, claims.Username)
	if err != nil {
		resp.FailWithMessage(utils.Translate(err), c)
		return
	}

	//如果修改了邮箱和用户名需要重新生成token
	if user.User.Email != nil || user.User.Email != nil {
		baseclaim := model.BaseClaims{
			Username: *user.User.Username,
			Email:    *user.User.Email,
		}
		genJwt := utils.NewJwt()
		tokenstr := genJwt.CreateClaims(baseclaim)
		//生成token
		token, err := genJwt.GenerateToken(&tokenstr)
		if err != nil {
			resp.FailWithMessage("token生成错误", c)
			return
		}

		//成功返回结果
		resp.OkWithData(token, c)
		return
	}
	resp.OkWithMessage("更新成功", c)

}
