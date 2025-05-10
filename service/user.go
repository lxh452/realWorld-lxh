package service

import (
	"realWorld/global"
	"realWorld/model"
	"realWorld/model/req"
	"realWorld/model/resp"
	"time"
)

var UserServiceApp = new(UserService)

type UserService struct {
}

func (u *UserService) Login(req req.UserAuthReq) (*resp.UserResp, error) {
	var user resp.UserResp
	err := global.DB.Where("email = ? and  password = ?", req.Email, req.Password).First(&user.User).Error
	if err != nil {
		return nil, global.ErrUserNotFound
	}
	return &user, nil
}

func (u *UserService) Register(req req.UserRegisterReq) (*resp.UserResp, error) {
	var user model.User
	var count int64
	// 检查用户是否存在
	err := global.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, global.ErrUserAlreadyExists
	}
	user.Username = req.Username
	user.Password = req.Passwd
	user.Email = req.Email
	err = global.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}
	var data resp.UserResp
	data.User.Bio = user.Bio
	data.User.Image = user.Image
	data.User.Email = user.Email
	return &data, nil
}

// 获取当前用户信息
func (u *UserService) GetUserInfo(req string) (*resp.UserResp, error) {
	var user resp.UserResp
	err := global.DB.Where("username = ?", req).First(&user.User).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 更改用户信息
func (u *UserService) ModifyUserInfo(req *req.ModifyUserInfo, reqName string) (*resp.UserResp, error) {
	var user resp.UserResp
	req.UpdatedAt = time.Now()
	result := global.DB.Debug().Where("username = ?", reqName).Updates(req)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
