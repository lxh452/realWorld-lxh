package model

// Tag 定义了数据库中的标签表结构
type Tag struct {
	Name string   `gorm:"column:name"` // 使用 gorm 标签指定数据库列名
	Tags []string `gorm:"-" json:"tags"`
}

// TableName 指定模型对应的数据库表名
func (t *Tag) TableName() string {
	return "tags"
}
