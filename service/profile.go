package service

import (
	"errors"
	"realWorld/global"
	"realWorld/model"
	"realWorld/model/resp"
)

// 获取用户信息包括是否自己关注了
// 需要关联表，用户表和关注表
var ProfileServiceApp = new(ProfileService)

type ProfileService struct{}

// 查看查找的用户信息
// 包含该用户信息和关联关注表查询的结果
func (profile *ProfileService) GetTagetUserInfo(follower *model.Follower) (*resp.ProfileResp, error) {
	var resp resp.ProfileResp
	err := global.DB.Table("users u").
		Select(`u.username,u.bio,u.image,CASE WHEN f.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS following`).
		Joins(`LEFT JOIN followers f ON f.user_id = ? AND f.follower_id = u.id AND f.deleted_at IS NULL`, follower.UserId). // 当前用户（关注发起者）
		Where("u.id = ?", follower.FollowerId).                                                                             // 目标用户
		Where("u.deleted_at IS NULL").
		Scan(&resp.Profile).Error
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

// 关注用户
func (profile *ProfileService) FollowUser(follower *model.Follower) (*resp.ProfileResp, error) {
	//查询是否已经关注
	var existingFollower model.Follower
	err := global.DB.Debug().Model(&model.Follower{}).
		Where("user_id = ? AND follower_id = ?", follower.UserId, follower.FollowerId).
		First(&existingFollower).Error
	if err == nil {
		return &resp.ProfileResp{}, errors.New("您已关注该用户")
	}
	//创建该关系表
	err = global.DB.Create(follower).Error
	if err != nil {
		return &resp.ProfileResp{}, errors.New("关注失败")
	}
	return profile.GetTagetUserInfo(follower)
}

// 取消关注
func (profile *ProfileService) UnFollowUser(follower *model.Follower) (*resp.ProfileResp, error) {
	//查询是否已经关注
	var existingFollower model.Follower
	err := global.DB.Debug().Model(&model.Follower{}).
		Where("user_id = ? AND follower_id = ?", follower.UserId, follower.FollowerId).
		First(&existingFollower).Error
	if err != nil {
		return &resp.ProfileResp{}, errors.New("您还没关注他")
	}
	//删除
	err = global.DB.Delete(&existingFollower).Error
	if err != nil {
		return &resp.ProfileResp{}, errors.New("取消关注失败")
	}
	return profile.GetTagetUserInfo(follower)
}
