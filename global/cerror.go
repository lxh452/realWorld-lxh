package global

import "errors"

var (
	// user
	ErrUserNotFound = errors.New("用户不存在")
	ErrPasswordIncorrect = errors.New("密码错误")
	ErrUserAlreadyExists = errors.New("用户已存在")
)

var (
	// role
	ErrRoleAlreadyExists = errors.New("角色已存在")
)