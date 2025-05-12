package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"realWorld/global"
	"realWorld/model/req"
	"realWorld/model/resp"
	"realWorld/service"
	"realWorld/utils"
	"strconv"
)

// 获取文章评论
func GetCommentFromArticleApi(c *gin.Context) {
	//获取文章摘要
	slug := c.Param("slug")
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
	err = faviorite.AddArticleToFaviorite(slug, claims.Id)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("添加文章到喜欢中"+err.Error(), zap.String("service", "AddArticleIntoFavoriteApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	resp.Ok(c)

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
	err = faviorite.DeleteArticleToFaviorite(slug, claims.Id)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("删除文章到喜欢中"+err.Error(), zap.String("service", "DeleteArticleFromFavoriteApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	resp.Ok(c)

}
