package config

// 管理redis的数据
type Redis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	DB       int    `json:"database"`
	Password string `json:"password"`
}
