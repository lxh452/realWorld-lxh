package service

import (
	"errors"
	"fmt"
	"realWorld/global"
	"realWorld/model"
	"realWorld/model/req"
	"realWorld/model/resp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var ArticleServiceApp = new(ArticleService)

type ArticleService struct{}

// 动态根据提供的字段进行修改查询方式
// 按条件获取文章列表
func (article *ArticleService) GetArticlesByConditions(tag string, authorId uint, favorited uint, reqid uint, limit string, offset string) ([]resp.ArticleResp, error) {
	var articles []resp.Articlegorm
	limitnum, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	offsetnum, err := strconv.Atoi(offset)
	if err != nil {
		return nil, err
	}

	query := global.DB.Model(&resp.Articlegorm{}).
		Select("articles.*, IFNULL(favorites.favorites_count, 0) AS favorites_count, IFNULL(user_favorites.favorited, 0) AS favorited").
		Joins("LEFT JOIN (SELECT article_id, COUNT(*) AS favorites_count FROM user_article_faviourite GROUP BY article_id) AS favorites ON articles.id = favorites.article_id").
		Joins("LEFT JOIN (SELECT article_id, COUNT(*) > 0 AS favorited FROM user_article_faviourite WHERE user_id = ? GROUP BY article_id) AS user_favorites ON articles.id = user_favorites.article_id", reqid).
		Order("articles.created_at DESC").
		Limit(limitnum).
		Offset(offsetnum)

	if tag != "" {
		query = query.Where("articles.tag_list LIKE ?", "%"+tag+"%")
	}
	if authorId > 0 {
		query = query.Where("articles.author_id = ?", authorId)
	}
	if favorited > 0 {
		query = query.Where("articles.id IN (SELECT article_id FROM user_article_faviourite WHERE user_id = ?)", favorited)
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
				Slug:           titleToSlug(articleinfo.Title),
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

	var articleinfo resp.Articlegorm
	// 使用 Raw 方法手动构建 SQL 查询
	err := global.DB.Model(&resp.Articlegorm{}).
		Select("articles.*, IFNULL(favorites.favorites_count, 0) AS favorites_count, IFNULL(user_favorites.favorited, 0) AS favorited").
		Joins("LEFT JOIN (SELECT article_id, COUNT(*) AS favorites_count FROM user_article_faviourite GROUP BY article_id) AS favorites ON articles.id = favorites.article_id").
		Joins("LEFT JOIN (SELECT article_id, COUNT(*) > 0 AS favorited FROM user_article_faviourite WHERE user_id = ? GROUP BY article_id) AS user_favorites ON articles.id = user_favorites.article_id", reqid).
		Where("articles.title = ?", slugToTitle(slug)).
		Scan(&articleinfo).Error
	if err != nil {
		return &resp.ArticleResp{}, err
	}
	// 调用私有方法获取作者信息
	info := &model.Follower{
		UserId:     reqid,
		FollowerId: articleinfo.AuthorId}
	authorinfo, err := getAuthorinfo(info)
	if err != nil {
		return nil, err
	}
	// 赋值
	data := resp.ArticleModel{
		Slug:           titleToSlug(articleinfo.Title),
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
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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
	result := global.DB.Where("title=? and author_id = ?", slugToTitle(slug), reqId).Updates(info)
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
	result := global.DB.Where("title= ? AND author_id = ?", slugToTitle(slug), reqId).First(&articleinfo)
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

// slugToTitle 将 slug 转换为 title
func slugToTitle(slug string) string {
	// 将连字符替换为空格
	title := strings.ReplaceAll(slug, "-", " ")

	// 将每个单词的首字母大写
	title = strings.Title(title)

	// 修正 strings.Title 的问题，它会将每个单词的后续字母都转换为小写
	title = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		return r
	}, title)

	return title
}

// titleToSlug 将 title 转换为 slug
func titleToSlug(title string) string {
	// 将所有字符转换为小写
	title = strings.ToLower(title)

	// 使用 strings.Builder 来构建最终的 slug
	var slug strings.Builder

	// 遍历每个字符
	for i, r := range title {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			// 如果是字母或数字，直接添加到 slug
			slug.WriteRune(r)
		} else if unicode.IsSpace(r) || r == '-' {
			// 如果是空格或连字符，转换为连字符
			if i > 0 && slug.Len() > 0 && slug.String()[slug.Len()-1] != '-' {
				slug.WriteRune('-')
			}
		}
		// 其他非字母数字字符忽略
	}

	// 移除开头和结尾的连字符
	slugStr := strings.Trim(slug.String(), "-")

	return slugStr
}
