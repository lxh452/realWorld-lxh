package req

import (
	"time"
)

type UserAuth struct {
	User UserAuthReq `json:"user"`
}

type UserRegister struct {
	User UserRegisterReq `json:"user"`
}

// 验证请求体
type UserAuthReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=12"`
}

// 注册请求体
type UserRegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=12"`
	Username string `json:"username" binding:"required,min=6,max=12"`
}

// ModifyUser 更改用户信息的请求体
type ModifyUser struct {
	User ModifyUserInfo `json:"user"`
}

// ModifyUserInfo 更改请求体
type ModifyUserInfo struct {
	Id        *uint     `json:"-"`
	Email     *string   `json:"email" binding:"omitempty,email"`
	Username  *string   `json:"username" binding:"omitempty,min=6,max=12"`
	Bio       *string   `json:"bio" binding:"omitempty"`
	Image     *string   `json:"image" binding:"omitempty,url"`
	Password  *string   `json:"password" binding:"omitempty,min=6,max=12"`
	UpdatedAt time.Time `json:"-"`
}

func (receiver ModifyUserInfo) TableName() string {
	return "users"
}
