package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"realWorld/global"
	"realWorld/model/req"
	"realWorld/model/resp"
	"realWorld/service"
	"realWorld/utils"
	"strconv"
	"time"
)

// 获取文章评论
func GetCommentFromArticleApi(c *gin.Context) {

	//获取文章摘要
	slug := c.Param("slug")

	//拼接key查询redis是否有该值
	cacheKey := fmt.Sprintf("Article_comments_%s", slug)
	cacheData, err := global.Redis.Get(context.TODO(), cacheKey).Bytes()
	if err == nil {
		var comment []resp.CommentResp
		if err = json.Unmarshal(cacheData, &comment); err != nil {
			resp.FailWithMessage("缓存数据格式错误", c)
			global.Logger.Warn("缓存数据格式错误"+err.Error(), zap.String("service", "GetCommentFromArticleApi"), zap.Int("port", global.CONFIG.Server.Port))
			return
		}
		resp.OkWithData(comment, c)
		return
	} else if err != redis.Nil {
		resp.FailWithMessage("缓存数据格式错误", c)
		global.Logger.Warn("缓存数据格式错误"+err.Error(), zap.String("service", "GetCommentFromArticleApi"), zap.Int("port", global.CONFIG.Server.Port))
		return

	}
	//绑定结构体

	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	comment := service.ArticleServiceApp
	articles, err := comment.GetCommentsFromArticle(slug, claims.Id)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("获取文章评论"+err.Error(), zap.String("service", "GetCommentFromArticleApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	cacheData, err = json.Marshal(&articles)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("缓存序列化失败"+err.Error(), zap.String("service", "GetCommentFromArticleApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	if _, err := global.Redis.Set(context.TODO(), cacheKey, cacheData, time.Minute*30).Result(); err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("缓存存储失败"+err.Error(), zap.String("service", "GetCommentFromArticleApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}

	resp.OkWithData(articles, c)

}

// 在文章添加评论
func AddcommentToArticleApi(c *gin.Context) {
	var commentReq req.CommentResp
	//获取文章摘要
	slug := c.Param("slug")
	//绑定结构体
	if err := c.ShouldBindJSON(&commentReq); err != nil {
		resp.FailWithMessage("绑定结构体失败"+err.Error(), c)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	comment := service.ArticleServiceApp
	article, err := comment.AddCommentToArticle(commentReq, claims.Id, slug)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("在文章添加评论"+err.Error(), zap.String("service", "AddcommentToArticleApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	resp.OkWithData(article, c)
}

// 删除文章评论(删除评论的id和评论id需为一致)
func DeleteCommentFromArticleApi(c *gin.Context) {
	slug := c.Param("slug")
	commentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	//获取用户数据
	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	comment := service.ArticleServiceApp
	err = comment.DeleteCommentFromArticle(slug, uint(commentId), claims.Id)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("删除文章评论"+err.Error(), zap.String("service", "DeleteCommentFromArticleApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	resp.Ok(c)
}

// 添加文章到喜欢中
func AddArticleIntoFavoriteApi(c *gin.Context) {
	slug := c.Param("slug")
	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	faviorite := service.ArticleServiceApp
	info, err := faviorite.AddArticleToFaviorite(slug, claims.Id)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("添加文章到喜欢中"+err.Error(), zap.String("service", "AddArticleIntoFavoriteApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	resp.OkWithData(info, c)

}

// 删除文章到喜欢中
func DeleteArticleFromFavoriteApi(c *gin.Context) {
	slug := c.Param("slug")
	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	faviorite := service.ArticleServiceApp
	info, err := faviorite.DeleteArticleToFaviorite(slug, claims.Id)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("删除文章到喜欢中"+err.Error(), zap.String("service", "DeleteArticleFromFavoriteApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	resp.OkWithData(info, c)

}
