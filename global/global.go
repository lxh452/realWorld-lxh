package global

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"realWorld/config"
)

var (
	CONFIG config.Config
	DB     *gorm.DB
	//Cron   *cron.Cron
	Redis  *redis.Client
	Logger *zap.Logger
)
