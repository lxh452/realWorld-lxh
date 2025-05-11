package resp

import "time"

//创建文章的接收体

type ArticleResp struct {
	Article ArticleModel `json:"article"`
}

// 创建文章的请求体
type ArticleModel struct {
	Slug           string      `json:"slug"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Body           string      `json:"body"`
	TagList        []string    `gorm:"type:text;serializer:json" json:"taglist"`
	CreatedAt      time.Time   `json:"createdAt"`
	UpdatedAt      time.Time   `json:"updatedAt"`
	Favorited      bool        `json:"favorited"`
	FavoritesCount uint        `json:"favoritescount"`
	Author         ProfileResp `json:"author"`
}

// 接收mysql数据
type Articlegorm struct {
	Id             uint      `json:"-"`
	AuthorId       uint      `json:"-"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	TagList        []string  `gorm:"type:text;serializer:json" json:"taglist"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `gorm:"->;column:favorited"`
	FavoritesCount uint      `gorm:"->;column:favorites_count"`
}

func (receiver Articlegorm) TableName() string {
	return "articles"
}

// ==================评论的结构体=====================
// 评论结构体
type CommentResp struct {
	Comment CommentModel `json:"comment"`
}

type CommentModel struct {
	ID        uint        `json:"id"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Body      string      `json:"body"`
	Author    ProfileResp `json:"author"`
}
