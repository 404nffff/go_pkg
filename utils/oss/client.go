package oss

import (
	"sync"

	"github.com/404nffff/go_pkg/oss"
	"github.com/404nffff/go_pkg/config"
)

var (
	clients = make(map[string]*oss.OSSClient)
	lock    sync.Mutex
)

// 初始化 OSS 客户端
func getClient(name string) (*oss.OSSClient, error) {
	lock.Lock()
	defer lock.Unlock()

	if client, ok := clients[name]; ok {
		return client, nil
	}

	//oss 前缀
	ossPrefix := "Oss." + name + "."

	ossConfig := oss.OssConfig{
		Endpoint:         config.ConfigYml.GetString(ossPrefix + "Endpoint"),
		AccessKeyID:      config.ConfigYml.GetString(ossPrefix + "AccessKeyId"),
		AccessKeySecret:  config.ConfigYml.GetString(ossPrefix + "AccessKeySecret"),
		BucketName:       config.ConfigYml.GetString(ossPrefix + "BucketName"),
		ConnectTimeout:   config.ConfigYml.GetInt(ossPrefix + "ConnectTimeout"),
		ReadWriteTimeout: config.ConfigYml.GetInt(ossPrefix + "ReadWriteTimeout"),
	}

	client, err := oss.NewOSSClient(ossConfig)
	if err != nil {
		return nil, err
	}

	clients[name] = client
	return client, nil
}
