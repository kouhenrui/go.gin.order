package repository

import mysql "go.gin.order/src/config/database"

var (
	db = mysql.NewSqlClient()
)

//func TableInit() {
//	db.AutoMigrate(acc, mer, cat)
//	log.Println("数据库表结构初始化成功")
//}
