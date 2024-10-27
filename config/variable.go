package config

import (
	"github.com/404nffff/go_pkg/ants"
	"github.com/404nffff/go_pkg/yml_config/ymlconfig_interf"

	"go.uber.org/zap"
)

// 初始化参数
var (
	BasePath string // 定义项目的根目录

	EventDestroyPrefix = "Destroy_" //  程序退出时需要销毁的事件前缀

	ConfigYml ymlconfig_interf.YmlConfigInterf // 全局配置文件指针

	Logs *zap.Logger // 全局日志指针

	Task ants.AntsInterface // 全局协程池指针
)
