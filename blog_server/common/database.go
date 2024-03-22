package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	user := "root"
	password := "123456"
	host := "localhost"
	port := "3306"
	database := "blog"
	timeout := "10s"
	driverName := "mysql"
	loc := "Asia%2FShanghai"
	args := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=true&loc=%s&timeout=%s",
		user,
		password,
		host,
		port,
		database,
		loc,
		timeout)
	// 连接数据库
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to open database: " + err.Error())
	}
	fmt.Println("成功连接数据库!")
	// 迁移数据表
	//db.AutoMigrate(&model.User{})
	//db.AutoMigrate(&model.Category{})
	//db.AutoMigrate(&model.Article{})
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
