package service

import (
	"gorm.io/gorm"
	pkgredis "wuzhispace.com/pkg/redis"
)

type BaseService struct {
	DB    *gorm.DB
	Redis *pkgredis.RedisClient
}
