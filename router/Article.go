package router

import (
	"github.com/gin-gonic/gin"
	"realWorld/api"
)

type ArticleGroup struct {
}

func (r *ArticleGroup) InitArticleRouters(engine *gin.Engine) {
	//获取标签
	Tags := engine.Group("/api/tags")
	{
		//获取所有标签返回标签列表
		Tags.GET("", api.GetTagsApi)
	}

	Article := engine.Group("/api/articles")
	//使用中间件
	Article.Use()
	{
		// 创建文章数量
		Article.POST("", api.CreateArticleApi)
		// 创建文章
		Article.GET("/:slug", api.GetArticleApi)
		// 更新文章
		Article.PUT("/:slug", api.UpdateArticleApi)
		//删除文章
		Article.DELETE("/:slug", api.DeleteArticleApi)

		//提要文章
		Article.GET("/feed", api.GetArticleFeedApi)
	}

	//无需使用中间件
	Article_NoAuth := engine.Group("/api/articles")
	{
		//按条件获取文章
		Article_NoAuth.GET("", api.GetArticlesApi)

	}

	//评论
	comments := engine.Group("/api/articles")
	//使用中间件
	comments.Use()
	{
		//获取该文章下的所有评论
		comments.GET("/:slug/comments", api.GetCommentFromArticleApi)
		//增加评论
		comments.POST("/:slug/comments", api.AddcommentToArticleApi)
		//删除评论
		comments.DELETE("/:slug/comments/:id", api.DeleteCommentFromArticleApi)

	}

	//喜欢
	faviorite := engine.Group("/api/articles")
	//使用中间件
	faviorite.Use()
	{
		faviorite.POST("/:slug/favorite", api.AddArticleIntoFavoriteApi)
		faviorite.DELETE("/:slug/favorite", api.DeleteArticleFromFavoriteApi)

	}

}
