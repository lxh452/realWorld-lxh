package service

import (
	"errors"
	"fmt"
	"realWorld/global"
	"realWorld/model"
	"realWorld/model/req"
	"realWorld/model/resp"
	"strconv"
	"time"
)

var ArticleServiceApp = new(ArticleService)

type ArticleService struct{}

// 动态根据提供的字段进行修改查询方式
// 按条件获取文章列表
func (article *ArticleService) GetArticlesByConditions(tag string, authorId uint, favorited string, limit string, offset string) ([]resp.ArticleResp, error) {
	var articles []resp.Articlegorm
	limitnum, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	offsetnum, err := strconv.Atoi(offset)
	if err != nil {
		return nil, err
	}

	query := global.DB.Model(&resp.Articlegorm{}).Debug().Order("created_at DESC").Limit(limitnum).Offset(offsetnum)

	if tag != "" {
		query = query.Where("tag_list LIKE ?", "%"+tag+"%")
	}
	if authorId > 0 {
		query = query.Where("author_id = ?", authorId)
	}
	if favorited != "" {
		query = query.Where("favorited_by LIKE ?", "%"+favorited+"%")
	}

	err = query.Find(&articles).Error
	if err != nil {
		return nil, err
	}

	// 将查询结果转换为嵌套了作者信息的结构体
	var articleResps []resp.ArticleResp
	for _, articleinfo := range articles {
		authorInfo, err := getAuthorinfo(&model.Follower{UserId: articleinfo.AuthorId, FollowerId: articleinfo.AuthorId})
		if err != nil {
			return nil, err
		}
		articleResp := resp.ArticleResp{
			Article: resp.ArticleModel{
				Slug:           articleinfo.Title,
				Title:          articleinfo.Title,
				Description:    articleinfo.Description,
				Body:           articleinfo.Body,
				TagList:        articleinfo.TagList,
				CreatedAt:      articleinfo.CreatedAt,
				UpdatedAt:      articleinfo.UpdatedAt,
				Favorited:      articleinfo.Favorited,
				FavoritesCount: articleinfo.FavoritesCount,
				Author:         *authorInfo,
			},
		}
		articleResps = append(articleResps, articleResp)
		fmt.Println(articleResp)
	}

	return articleResps, nil
}

// 获取文章信息
func (article *ArticleService) GetArticleInfo(slug string, reqid uint) (*resp.ArticleResp, error) {
	fmt.Println(slug, reqid)
	var articleinfo resp.Articlegorm
	err := global.DB.Model(&resp.Articlegorm{}).Where("title = ?", slug).Scan(&articleinfo).Error
	if err != nil {
		return &resp.ArticleResp{}, err
	}
	//调用私有方法获取作者信息
	info := &model.Follower{
		UserId:     reqid,
		FollowerId: articleinfo.AuthorId}
	authorinfo, err := getAuthorinfo(info)
	if err != nil {
		return nil, err
	}
	//赋值
	data := resp.ArticleModel{
		Slug:           articleinfo.Title,
		Title:          articleinfo.Title,
		Description:    articleinfo.Description,
		Body:           articleinfo.Body,
		TagList:        articleinfo.TagList,
		CreatedAt:      articleinfo.CreatedAt,
		UpdatedAt:      articleinfo.UpdatedAt,
		Favorited:      articleinfo.Favorited,
		FavoritesCount: articleinfo.FavoritesCount,
		Author:         *authorinfo,
	}
	resp := &resp.ArticleResp{Article: data}
	fmt.Println(resp)
	return resp, nil
}

// 创建文章
func (article *ArticleService) CreateArticle(req *req.CreateArticleReq, authorid uint) (*resp.ArticleResp, error) {
	//将请求体跟获取的用户id进行绑定
	data := model.Article{
		Title:       req.Title,
		Body:        req.Body,
		TagList:     req.TagList,
		AuthorID:    authorid,
		Description: req.Description,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
	}
	//处理错误
	result := global.DB.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	// 获取文章信息
	resp, err := article.GetArticleInfo(req.Title, authorid)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 更新文章信息
func (article *ArticleService) UpdateArticle(info *req.ModifyArticleInfo, slug string, reqId uint) (*resp.ArticleResp, error) {
	info.UpdatedAt = time.Now()
	result := global.DB.Where("title=? and author_id = ?", slug, reqId).Updates(info)
	if result.Error != nil {
		return &resp.ArticleResp{}, result.Error
	}

	//如果标题更新了，根据新的标题查询
	if info.Title == nil {
		return article.GetArticleInfo(slug, reqId)
	}
	return article.GetArticleInfo(*info.Title, reqId)
}

// 删除文章
func (article *ArticleService) DeleteArticle(slug string, reqId uint) error {
	var articleinfo model.Article
	result := global.DB.Where("title= ? AND author_id = ?", slug, reqId).First(&articleinfo)
	if result.Error != nil {
		return errors.New("文章不存在或你没有权限删除")
	}
	if result.RowsAffected == 0 {
		return errors.New("文章不存在或你没有权限删除")
	}
	err := global.DB.Delete(&articleinfo).Error
	if err != nil {
		return err
	}
	return nil
}

// 获取作者信息
// 私有方法获取作者个人信息，调用profile接口
func getAuthorinfo(req *model.Follower) (resp *resp.ProfileResp, err error) {
	profile := ProfileService{}
	Author, err := profile.GetTagetUserInfo(req)
	if err != nil {
		return resp, errors.New("获取个人信息失败")
	}
	return Author, nil
}
