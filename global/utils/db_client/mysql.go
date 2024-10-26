package db_client

import (
	"github.com/404nffff/go_pkg/pkg/mysql"

	"gorm.io/gorm"
)

func MysqlLocal() *gorm.DB {

	return mysql.NewClient("Local")
}
