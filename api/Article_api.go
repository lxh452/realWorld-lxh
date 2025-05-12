package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"realWorld/global"
	"realWorld/model/req"
	"realWorld/model/resp"
	"realWorld/service"
	"realWorld/utils"
)

// 获取所有标签
func GetTagsApi(c *gin.Context) {
	Tag := service.TagsServiceApp
	tags, err := Tag.GetAllTags()
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("获取标签列表"+err.Error(), zap.String("service", "GetTagsApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	resp.OkWithData(tags, c)
}

// 按条件获取文章列表
func GetArticlesApi(c *gin.Context) {
	limit := c.DefaultQuery("limit", "20")  // 限制文章数量，默认为20
	offset := c.DefaultQuery("offset", "0") // 偏移/跳过文章数量，默认为0
	tag := c.Query("tag")                   // 按标签过滤
	author := c.Query("author")             // 按作者筛选
	favorited := c.Query("favorited")       // 指定用户的收藏
	//根据用户名查找id
	//获取自身id
	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	articles := service.ArticleServiceApp
	conditions, err := articles.GetArticlesByConditions(tag, getTargetId(author), getTargetId(favorited), claims.Id, limit, offset)
	if err != nil {
		return
	}
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("按条件获取文章列表"+err.Error(), zap.String("service", "GetArticlesApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	resp.OkWithData(conditions, c)
}

// 获取单个文章
func GetArticleApi(c *gin.Context) {
	slug := c.Param("slug")
	//从token中获取数据
	claims, err := utils.GetClaims(c)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		return
	}
	article := service.ArticleServiceApp
	info, err := article.GetArticleInfo(slug, claims.Id)
	if err != nil {
		resp.FailWithMessage(err.Error(), c)
		//写入日志
		global.Logger.Warn("获取单个文章"+err.Error(), zap.String("service", "GetArticleApi"), zap.Int("port", global.CONFIG.Server.Port))
		return
	}
	resp.OkWithData(info, c)
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
