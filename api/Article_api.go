package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"realWorld/global"
	"realWorld/model/req"
	"realWorld/model/resp"
	"realWorld/service"
	"realWorld/utils"
	"time"
)

// 获取所有标签
func GetTagsApi(c *gin.Context) {
	// 尝试从 Redis 缓存中获取标签
	cacheData, err := global.Redis.Get(context.TODO(), "tags").Bytes() // 获取字节数据
	if err == nil {
		// 缓存命中，反序列化缓存数据
		var tags resp.TagResp
		if err := json.Unmarshal(cacheData, &tags); err != nil {
			resp.FailWithMessage("缓存数据格式错误", c)
			global.Logger.Warn("缓存数据格式错误", zap.Error(err), zap.String("service", "GetTagsApi"), zap.Int("port", global.CONFIG.Server.Port))
			return
		}
		resp.OkWithData(tags, c)
		return
	} else if err != redis.Nil {
		// Redis 错误，记录日志并返回错误
		resp.FailWithMessage("缓存获取失败", c)
		global.Logger.Warn("缓存获取失败", zap.Error(err), zap.String("service", "GetTagsApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}

	// 缓存未命中，从数据库中获取标签
	tagService := service.TagsServiceApp
	tags, err := tagService.GetAllTags()
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		global.Logger.Warn("获取标签列表失败", zap.Error(err), zap.String("service", "GetTagsApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}

	// 将获取到的标签存储到 Redis 缓存中
	cacheData, err = json.Marshal(tags)
	if err != nil {
		resp.FailWithMessage("缓存数据序列化失败", c)
		global.Logger.Warn("缓存数据序列化失败", zap.Error(err), zap.String("service", "GetTagsApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	if _, err := global.Redis.Set(context.TODO(), "tags", cacheData, time.Minute*30).Result(); err != nil {
		resp.FailWithMessage("缓存存储失败", c)
		global.Logger.Warn("缓存存储失败", zap.Error(err), zap.String("service", "GetTagsApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}

	// 返回标签数据
	resp.OkWithData(tags, c)
}

// 按条件获取文章列表
func GetArticlesApi(c *gin.Context) {
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")
	tag := c.Query("tag")
	author := c.Query("author")
	favorited := c.Query("favorited")

	cacheKey := fmt.Sprintf("articles:tag:%s:author:%s:favorited:%s:limit:%s:offset:%s", tag, author, favorited, limit, offset)
	cacheData, err := global.Redis.Get(context.TODO(), cacheKey).Bytes()
	if err == nil {
		// 缓存命中，反序列化缓存数据
		var articles []resp.ArticleResp
		if err := json.Unmarshal(cacheData, &articles); err != nil {
			resp.FailWithMessage("缓存数据格式错误", c)
			global.Logger.Warn("缓存数据格式错误", zap.Error(err), zap.String("service", "GetArticlesApi"), zap.Int("port", global.CONFIG.Server.Port))
			return
		}
		resp.OkWithData(articles, c)
		return
	} else if err != redis.Nil {
		// Redis 错误，记录日志并返回错误
		resp.FailWithMessage("缓存获取失败", c)
		global.Logger.Warn("缓存获取失败", zap.Error(err), zap.String("service", "GetArticlesApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}

	// 缓存未命中，从数据库中获取文章列表
	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	articles, err := service.ArticleServiceApp.GetArticlesByConditions(tag, getTargetId(author), getTargetId(favorited), claims.Id, limit, offset)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		global.Logger.Warn("按条件获取文章列表失败", zap.Error(err), zap.String("service", "GetArticlesApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}

	// 将获取到的文章列表存储到 Redis 缓存中
	cacheData, err = json.Marshal(articles)
	if err != nil {
		resp.FailWithMessage("缓存数据序列化失败", c)
		global.Logger.Warn("缓存数据序列化失败", zap.Error(err), zap.String("service", "GetArticlesApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	if _, err := global.Redis.Set(context.TODO(), cacheKey, cacheData, time.Minute*30).Result(); err != nil {
		resp.FailWithMessage("缓存存储失败", c)
		global.Logger.Warn("缓存存储失败", zap.Error(err), zap.String("service", "GetArticlesApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}

	// 返回文章列表
	resp.OkWithData(articles, c)
}

// 获取单个文章
func GetArticleApi(c *gin.Context) {
	slug := c.Param("slug")
	cacheKey := fmt.Sprintf("article:%s", slug)
	cacheData, err := global.Redis.Get(context.TODO(), cacheKey).Bytes()
	if err == nil {
		// 缓存命中，反序列化缓存数据
		var article resp.ArticleResp
		if err := json.Unmarshal(cacheData, &article); err != nil {
			resp.FailWithMessage("缓存数据格式错误", c)
			global.Logger.Warn("缓存数据格式错误", zap.Error(err), zap.String("service", "GetArticleApi"), zap.Int("port", global.CONFIG.Server.Port))
			return
		}
		resp.OkWithData(article, c)
		return
	} else if err != redis.Nil {
		// Redis 错误，记录日志并返回错误
		resp.FailWithMessage("缓存获取失败", c)
		global.Logger.Warn("缓存获取失败", zap.Error(err), zap.String("service", "GetArticleApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}

	// 缓存未命中，从数据库中获取文章
	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	article, err := service.ArticleServiceApp.GetArticleInfo(slug, claims.Id)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		global.Logger.Warn("获取单个文章失败", zap.Error(err), zap.String("service", "GetArticleApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}

	// 将获取到的文章存储到 Redis 缓存中
	cacheData, err = json.Marshal(article)
	if err != nil {
		resp.FailWithMessage("缓存数据序列化失败", c)
		global.Logger.Warn("缓存数据序列化失败", zap.Error(err), zap.String("service", "GetArticleApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	if _, err := global.Redis.Set(context.TODO(), cacheKey, cacheData, time.Minute*30).Result(); err != nil {
		resp.FailWithMessage("缓存存储失败", c)
		global.Logger.Warn("缓存存储失败", zap.Error(err), zap.String("service", "GetArticleApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}

	// 返回文章数据
	resp.OkWithData(article, c)
}

// 创建文章
func CreateArticleApi(c *gin.Context) {
	//定义一个请求体接收数据
	var articlereq req.CreateArticle
	err := c.ShouldBindJSON(&articlereq)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	//	2.进行监管，敏感性检测，存放到临时表
	//	3.在token中拿取id信息，
	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	//获取个人资料调用profile的接口
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	//	4.将文件保存到持久层
	article := service.ArticleServiceApp
	fmt.Println(claims)
	createArticle, err := article.CreateArticle(&articlereq.Article, claims.Id)
	fmt.Println(createArticle)
	//返回错误
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("创建文章"+err.Error(), zap.String("service", "CreateArticleApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	//返回成功结果
	resp.OkWithData(createArticle, c)

}

// 更新文章
func UpdateArticleApi(c *gin.Context) {
	slug := c.Param("slug")
	//从token中获取数据
	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	//绑定结构体
	var articleReq req.UpdateArticle
	if err = c.ShouldBindJSON(&articleReq); err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	//校验
	validate := validator.New()
	if err = validate.Struct(articleReq); err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	article := service.ArticleServiceApp
	updateArticle, err := article.UpdateArticle(&articleReq.Article, slug, claims.Id)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("更新文章"+err.Error(), zap.String("service", "UpdateArticleApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	resp.OkWithData(updateArticle, c)

}

//删除文章

func DeleteArticleApi(c *gin.Context) {
	slug := c.Param("slug")
	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	article := service.ArticleService{}
	err = article.DeleteArticle(slug, claims.Id)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("删除文章"+err.Error(), zap.String("service", "DeleteArticleApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	resp.OkWithMessage("删除成功", c)
}

// 获取提要文章
func GetArticleFeedApi(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	limit := c.DefaultQuery("limit", "20")  // 限制文章数量，默认为20
	offset := c.DefaultQuery("offset", "0") // 偏移/跳过文章数量，默认为0
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	articles := service.ArticleServiceApp
	followedArticles, err := articles.GetFollowedArticles(claims.Id, limit, offset)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	resp.OkWithData(followedArticles, c)
}
