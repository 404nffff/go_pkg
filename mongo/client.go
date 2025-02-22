package mongo

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/404nffff/go_pkg/event_manage"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// 全局 sync.Map 变量
var (
	dbs sync.Map
)

// NewClient 初始化 MongoDB 客户端，并支持多个数据库连接
func NewClient(configName string) *mongo.Database {
	// 使用 LoadOrStore 确保在并发环境中只初始化一次数据库连接
	db, loaded := dbs.LoadOrStore(configName, func() interface{} {
		return createMongoClient(configName)
	})

	database := db.(*mongo.Database)

	if loaded {
		// 检查连接是否有效
		client := database.Client()
		if !isValidConnection(client) {
			log.Printf("MongoDB 连接丢失，正在重新连接: %s", configName)
			database = createMongoClient(configName)
			dbs.Store(configName, database)
		}
	}
	return database
}

// isValidConnection 检查 MongoDB 连接是否有效
func isValidConnection(client *mongo.Client) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := client.Ping(ctx, readpref.Primary())
	return err == nil
}

// createMongoClient 创建新的 MongoDB 客户端
func createMongoClient(configName string) *mongo.Database {

	// 加载配置
	dbConfig := loadConfig(configName)

	// 判断是否为空
	if dbConfig.URI == "" {
		panic(fmt.Sprintf("Failed to get MongoDB config: %s", configName))
	}

	// 设置客户端连接选项
	clientOptions := options.Client().
		ApplyURI(dbConfig.URI + dbConfig.Database).
		SetMaxPoolSize(dbConfig.MaxPoolSize).
		SetMinPoolSize(dbConfig.MinPoolSize).
		SetMonitor(NewMonitor())

	// 连接到 MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
	}

	// 检查连接
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(fmt.Sprintf("Failed to ping MongoDB: %v", err))
	}

	log.Printf("Connected to MongoDB successfully, database: %s", dbConfig.Database)

	// 注册销毁事件
	eventManageFactory := event_manage.CreateEventManageFactory()
	eventName := dbConfig.EventDestroyPrefix
	if _, exists := eventManageFactory.Get(eventName); !exists {
		eventManageFactory.Set(eventName, func(args ...interface{}) {
			CloseMongo(client, dbConfig.Database)
			log.Printf("Destroying MongoDB connection for %s", dbConfig.Database)
		})
	}

	return client.Database(dbConfig.Database)
}

// GetCollection 获取指定数据库的集合
func GetCollection(dbName string, collection string) *mongo.Collection {
	db, exists := dbs.Load(dbName)
	if exists {
		return db.(*mongo.Database).Collection(collection)
	}
	return nil
}

// FindOne 查找单个文档
func FindOne(dbName string, collection string, filter interface{}) *mongo.SingleResult {
	col := GetCollection(dbName, collection)
	if col == nil {
		return nil
	}
	return col.FindOne(context.TODO(), filter)
}

// CloseMongo 关闭 MongoDB 客户端
func CloseMongo(client *mongo.Client, name string) {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
		log.Printf("Disconnected from MongoDB successfully, client: %s", name)
	}
}
