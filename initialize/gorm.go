package initialize

import (
	"fmt"
	"realWorld/global"
	"realWorld/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Article{},
		&model.Comment{},
		&model.Follower{},
	)
}

func MustLoadGorm() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", global.CONFIG.MySQL.Username, global.CONFIG.MySQL.Password,
		global.CONFIG.MySQL.Host, global.CONFIG.MySQL.Port, global.CONFIG.MySQL.Database)
	fmt.Println("我的mysql", dsn)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	global.DB = db
}
