package db_client

import (
	"github.com/404nffff/go_pkg/pkg/memcached"

	"github.com/bradfitz/gomemcache/memcache"
)

// 本地连接 Memcached
func MemLocal() *memcache.Client {
	return memcached.NewClient("Local")
}
