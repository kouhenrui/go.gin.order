package database

import (
	"fmt"
	"go.gin.order/src/internal/pojo"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"time"
)

var (
	ClickhouseClient *gorm.DB
)

func InitClickhouse(conf *pojo.ClickConf) {
	dsn := fmt.Sprintf("clickhouse://gorm:gorm@%s:%s/%sdial_timeout=10s&read_timeout=20s", conf.Host, conf.Port, conf.Database)
	//dsn := "clickhouse://gorm:gorm@localhost:9942/gorm?dial_timeout=10s&read_timeout=20s"
	ClickhouseClient, err = gorm.Open(clickhouse.Open(dsn), &gorm.Config{})

	sqldb, _ := ClickhouseClient.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqldb.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqldb.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqldb.SetConnMaxLifetime(time.Hour)
}
