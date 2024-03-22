package main

import (
	"blog_server/common"
	"blog_server/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func main() {
	db := common.InitDB()
	defer db.Close()
	r := gin.Default()
	//配置静态文件路径
	r.StaticFS("/images", http.Dir("./static/images"))
	routes.CollectRoutes(r)
	panic(r.Run(":8080"))
}
