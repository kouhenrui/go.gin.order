package database

import (
	"fmt"
	"go.gin.order/src/internal/pojo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

var (
	PostgreClient *gorm.DB
)

func InitPostgre(postgre *pojo.PostGreConf) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", postgre.Host, postgre.User, postgre.Password, postgre.Database, postgre.Port)
	log.Println(dsn)
	//dsn := "host=localhost user=youruser password=yourpassword dbname=yourdbname port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	PostgreClient, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:         "",   //无前缀
			SingularTable:       true, //表名无复数
			NameReplacer:        nil,
			NoLowerCase:         false,
			IdentifierMaxLength: 0,
		},
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	sqldb, _ := PostgreClient.DB()
	// 设置数据库连接池参数
	sqldb.SetMaxOpenConns(100)          // 设置最大打开连接数
	sqldb.SetMaxIdleConns(10)           // 设置最大空闲连接数
	sqldb.SetConnMaxLifetime(time.Hour) // 设置连接的最大生命周期
	InitPostgresql()
	log.Println("postgresql连接成功，表初始化成功")
}

var (
	car   pojo.Cart
	catit pojo.CartItem
	ord   pojo.Order
	ordit pojo.OrderItem

	apl pojo.Approval
	apr pojo.Approver
	apn pojo.ApprovalAction
)

func InitPostgresql() {
	PostgreClient.AutoMigrate(&car, &ord, &catit, &ordit, &apl, &apr, &apn)
}
