package model

import "gorm.io/gorm"

// 关注用户的中间表 多对多
// 一个用户可以关注多个用户，多个用户也可以关注自己

// Follower 表示用户关注关系
type Follower struct {
	gorm.Model
	UserId     uint `gorm:"uniqueIndex:idx_user_follower"`
	FollowerId uint `gorm:"uniqueIndex:idx_user_follower"`
}
