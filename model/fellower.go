package model

import "gorm.io/gorm"

// 关注用户的中间表 多对多
// 一个用户可以关注多个用户，多个用户也可以关注自己

type Follower struct {
	gorm.Model
	UserId     uint
	FollowerId uint
}
