package model

import "gorm.io/gorm"

// User 表示用户信息
type User struct {
	gorm.Model
	Email    string `gorm:"unique;index:idx_email"`    // 唯一索引
	Username string `gorm:"unique;index:idx_username"` // 唯一索引
	Password string
	Bio      string
	Image    *string
	Article  []Article `gorm:"many2many:user_article_faviourite"`
	//Follows  []User    `gorm:"many2many:user_follows;joinForeignKey:UserID;joinReferences:FollowID"`
}
