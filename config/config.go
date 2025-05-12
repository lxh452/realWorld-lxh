package config

type Config struct {
	MySQL  MySQL
	Server Server
	Jwt    Jwt
	Logs   Logs
	Redis  Redis
}
