package db_client

import (
	pkgRedis "github.com/404nffff/go_pkg/pkg/redis"

	"github.com/go-redis/redis/v8"
)

func RedisLocal() *redis.Client {

	return pkgRedis.NewClient("Local")
}
