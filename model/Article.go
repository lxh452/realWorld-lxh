package model

import (
	"time"
)

// Article 表示文章主体
type Article struct {
	AuthorID    uint
	Author      User      `gorm:"foreignkey:AuthorID"`
	ID          uint      `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
	Title       string    `gorm:"not_null" json:"title"`
	Description string    `gorm:"not_null" json:"description"`
	Body        string    `gorm:"not_null" json:"body"`
	TagList     []string  `gorm:"type:text;serializer:json" json:"taglist"`
	Comments    []Comment `gorm:"foreignkey:ArticleID" json:"comments"`
	User        []User    `gorm:"many2many:user_article_faviourite"`
}

// Comment 表示评论
type Comment struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
	ArticleID uint
	Article   Article `gorm:"foreignkey:ArticleID"`
	Body      string  `gorm:"not_null"`
	AuthorID  uint
	Author    User `gorm:"foreignkey:AuthorID"`
}
