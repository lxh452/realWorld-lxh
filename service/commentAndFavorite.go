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

//service层的article共用一个结构体，创建方法

// 增加评论到评论表中
// 1.接收评论请求体，获取当前用户id，获取文章id，绑定到评论表
// 查询优化，只查询作者id
func (comment ArticleService) AddCommentToArticle(req req.CommentResp, userId uint, slug string) (*resp.CommentResp, error) {
	// 1.查询文章ID
	var articleId uint
	if err := global.DB.Table("articles").Select("id").Where("title = ?", slug).Scan(&articleId).Error; err != nil {
		return &resp.CommentResp{}, err
	}
	// 2.创建评论插入到数据库
	commentData := model.Comment{
		ArticleID: articleId,
		AuthorID:  userId,
		Body:      req.Comment.Body,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// 3.插入评论到数据库
	if err := global.DB.Create(&commentData).Error; err != nil {
		return &resp.CommentResp{}, err
	}
	//返回回应结构体
	fmt.Println(userId)
	return comment.GetCommentFromArticle(slug, userId)
}

// 获取单条评论
func (comment ArticleService) GetCommentFromArticle(slug string, userId uint) (*resp.CommentResp, error) {
	var commentData model.Comment
	var articleId uint
	//获取文章中作者的id
	if err := global.DB.Table("articles").Select("author_id").Where("title = ?", slug).Scan(&articleId).Error; err != nil {
		return &resp.CommentResp{}, err
	}
	//2.根据文章id和评论人的个人id
	err := global.DB.Table("comments").Where("article_id = ? and author_id = ?", articleId, userId).Scan(&commentData).Error
	if err != nil {
		return &resp.CommentResp{}, err
	}
	//	获取评论者和用户自身的关系
	relationship := &model.Follower{
		UserId:     userId,
		FollowerId: commentData.AuthorID,
	}
	authorInfo, err := getAuthorinfo(relationship)
	if err != nil {
		return nil, err
	}
	commentModel := resp.CommentModel{
		ID:        commentData.ID,
		CreatedAt: commentData.CreatedAt,
		UpdatedAt: commentData.UpdatedAt,
		Body:      commentData.Body,
		Author:    *authorInfo,
	}
	return &resp.CommentResp{Comment: commentModel}, nil
}

// GetCommentsFromArticle 获取文章的多条评论
func (comment ArticleService) GetCommentsFromArticle(slug string, userId uint) ([]resp.CommentResp, error) {
	var comments []model.Comment
	var articleId uint
	if err := global.DB.Table("articles").Select("id").Where("title = ?", slug).Scan(&articleId).Error; err != nil {
		return nil, err
	}
	err := global.DB.Table("comments").Where("article_id = ?", articleId).Find(&comments).Error
	if err != nil {
		return nil, err
	}

	var commentResps []resp.CommentResp
	for _, commentData := range comments {
		relationship := &model.Follower{
			UserId:     userId,
			FollowerId: commentData.AuthorID,
		}
		authorinfo, err := getAuthorinfo(relationship)
		if err != nil {
			return nil, err
		}
		model := resp.CommentModel{
			ID:        commentData.ID,
			CreatedAt: commentData.CreatedAt,
			UpdatedAt: commentData.UpdatedAt,
			Body:      commentData.Body,
			Author:    *authorinfo,
		}
		commentResps = append(commentResps, resp.CommentResp{Comment: model})
	}
	return commentResps, nil
}

// DeleteCommentFromArticle 删除评论
func (comment ArticleService) DeleteCommentFromArticle(slug string, commentId uint, userId uint) error {
	var articleId uint
	if err := global.DB.Table("articles").Select("id").Where("title = ?", slug).Scan(&articleId).Error; err != nil {
		return err
	}

	var commentData model.Comment
	if err := global.DB.Where("id = ? AND article_id = ? AND author_id = ?", commentId, articleId, userId).First(&commentData).Error; err != nil {
		return errors.New("评论不存在或你没有权限删除")
	}

	if err := global.DB.Delete(&commentData).Error; err != nil {
		return err
	}
	return nil
}

// 创建喜欢文章
func (faviorite ArticleService) AddArticleToFaviorite(slug string, userId uint) error {
	var articleId uint
	if err := global.DB.Table("articles").Select("id").Where("title = ?", slug).Scan(&articleId).Error; err != nil {
		return err
	}
	relation := &req.Faviorite{
		UserId:    userId,
		ArticleId: articleId,
	}
	err := global.DB.Table("user_article_faviourite").Create(&relation).Error
	if err != nil {
		return err
	}
	return nil
}

// 删除喜欢文章
func (faviorite ArticleService) DeleteArticleToFaviorite(slug string, userId uint) error {
	var articleId uint
	if err := global.DB.Table("articles").Select("id").Where("title = ?", slug).Scan(&articleId).Error; err != nil {
		return err
	}

	// 明确指定删除条件
	err := global.DB.Table("user_article_faviourite").Where("user_id = ? AND article_id = ?", userId, articleId).Delete(&req.Faviorite{}).Error
	if err != nil {
		return err
	}
	return nil
}

// 获取关注用户的文章列表
func (article *ArticleService) GetFollowedArticles(reqid uint, limit string, offset string) ([]resp.ArticleResp, error) {
	var articles []resp.Articlegorm
	limitnum, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	offsetnum, err := strconv.Atoi(offset)
	if err != nil {
		return nil, err
	}

	// 查询关注用户的文章
	query := global.DB.Model(&resp.Articlegorm{}).
		Select("articles.*, IFNULL(favorites.favorites_count, 0) AS favorites_count, IFNULL(user_favorites.favorited, 0) AS favorited").
		Joins("LEFT JOIN (SELECT article_id, COUNT(*) AS favorites_count FROM user_article_faviourite GROUP BY article_id) AS favorites ON articles.id = favorites.article_id").
		Joins("LEFT JOIN (SELECT article_id, COUNT(*) > 0 AS favorited FROM user_article_faviourite WHERE user_id = ? GROUP BY article_id) AS user_favorites ON articles.id = user_favorites.article_id", reqid).
		Joins("INNER JOIN followers f ON f.follower_id = articles.author_id").
		Where("f.user_id = ? AND f.deleted_at IS NULL", reqid).
		Order("articles.created_at DESC").
		Limit(limitnum).
		Offset(offsetnum)

	err = query.Find(&articles).Error
	if err != nil {
		return nil, err
	}

	// 将查询结果转换为嵌套了作者信息的结构体
	var articleResps []resp.ArticleResp
	for _, articleinfo := range articles {
		authorInfo, err := getAuthorinfo(&model.Follower{UserId: reqid, FollowerId: articleinfo.AuthorId})
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
	}

	return articleResps, nil
}
