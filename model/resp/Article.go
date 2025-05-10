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
	FavoritesCount uint64      `json:"favoritescount"`
	Author         ProfileResp `json:"author"`
}

// 接收mysql数据
type Articlegorm struct {
	AuthorId       uint      `json:"-"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	TagList        []string  `gorm:"type:text;serializer:json" json:"taglist"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount uint64    `json:"favoritescount"`
}

func (receiver Articlegorm) TableName() string {
	return "articles"
}
