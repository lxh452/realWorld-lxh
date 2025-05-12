package service

import (
	"realWorld/global"
	"realWorld/model"
	"realWorld/model/resp"
)

var TagsServiceApp = new(TagsService)

type TagsService struct{}

// GetAllTags 获取所有标签
func (tagsService *TagsService) GetAllTags() (*resp.TagResp, error) {
	var tagResp []model.Tag
	result := global.DB.Find(&tagResp)
	if result.Error != nil {
		return nil, result.Error
	}

	tagSlice := make([]string, 0) // 初始化切片，长度为0，因为我们不知道有多少标签
	for _, tag := range tagResp {
		if tag.Name != "" { // 检查标签是否为空
			tagSlice = append(tagSlice, tag.Name)
		}
	}

	return &resp.TagResp{Tags: tagSlice}, nil
}
