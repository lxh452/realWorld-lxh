package router

import (
	"github.com/gin-gonic/gin"
	"realWorld/api"
)

type ArticleGroup struct {
}

func (r *ArticleGroup) InitArticleRouters(engine *gin.Engine) {
	Article := engine.Group("/api/articles")
	//使用中间件
	Article.Use()
	{
		// 创建文章数量
		Article.POST("", api.CreateArticle)
		// 创建文章
		Article.GET("/:slug", api.GetArticle)
		// 更新文章
		Article.PUT("/:slug", api.UpdateArticle)
		//删除文章
		Article.DELETE("/:slug", api.DeleteArticle)
	}
	Article_NoAuth := engine.Group("/api/articles")
	{
		//按条件获取文章
		Article_NoAuth.GET("", api.GetArticles)

	}

	//评论
	comments := engine.Group("/api/articles")
	//使用中间件
	comments.Use()
	{
		//获取该文章下的所有评论
		comments.GET("/:slug/comments", api.GetCommentFromArticle)
		//增加评论
		comments.POST("/:slug/comments", api.AddcommentToArticle)
		//删除评论
		comments.DELETE("/:slug/comments/:id", api.DeleteCommentFromArticle)

	}
	//喜欢
	faviorite := engine.Group("/api/articles")
	//使用中间件
	faviorite.Use()
	{
		faviorite.POST("/:slug/favorite", api.AddArticleIntoFavorite)
		faviorite.DELETE("/:slug/favorite", api.DeleteArticleFromFavorite)

	}

}
