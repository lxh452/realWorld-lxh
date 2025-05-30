package config

// 用来管理MySQL的配置
type MySQL struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Charset  string `json:"charset"`
	Config   string `json:"config"`
}
