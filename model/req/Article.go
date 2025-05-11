package req

import (
	"time"
)

type CreateArticle struct {
	Article CreateArticleReq `json:"article"`
}
type CreateArticleReq struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList"`
}
type UpdateArticle struct {
	Article ModifyArticleInfo `json:"article"`
}
type ModifyArticleInfo struct {
	Title       *string   `json:"title" binding:"omitempty"`
	Description *string   `json:"description" binding:"omitempty"`
	Body        *string   `json:"body" binding:"omitempty"`
	UpdatedAt   time.Time `json:"-"`
}

func (u ModifyArticleInfo) TableName() string {
	return "articles"
}

// 评论请求体
type CommentResp struct {
	Comment commentModel `json:"comment"`
}
type commentModel struct {
	Body string `json:"body"`
}

type Faviorite struct {
	ArticleId uint `json:"article_id"`
	UserId    uint `json:"user_id"`
}
