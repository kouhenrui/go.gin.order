package database

import (
	"context"
	"fmt"
	"go.gin.order/src/internal/pojo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	MongoDatabase *mongo.Database
)

func InitMongo(conf pojo.MongoConf) {
	clientoptions := options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%s", conf.Ip, conf.Port)).
		SetMaxPoolSize(100).                       // 设置连接池的最大连接数
		SetMinPoolSize(10).                        // 设置连接池的最小连接数
		SetMaxConnIdleTime(30 * time.Second).      // 设置连接的最大空闲时间
		SetConnectTimeout(10 * time.Second).       // 设置连接超时时间
		SetServerSelectionTimeout(5 * time.Second) // 设置服务器选择超时时间
	client, err := mongo.NewClient(clientoptions)
	if err != nil {
		if err != nil {
			panic("连接mongo数据库失败, error=" + err.Error())
		}
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB!")
	MongoDatabase = client.Database(conf.Database)
}
