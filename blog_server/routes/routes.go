package routes

import (
	"blog_server/controller"
	"blog_server/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoutes(r *gin.Engine) *gin.Engine {
	//允许跨域访问
	r.Use(middleware.CORSMiddleware())
	//注册
	r.POST("/register", controller.Register)
	//登录
	r.POST("/login", controller.Login)
	//上传头像
	r.POST("/upload", controller.Upload)
	r.POST("/upload/rich_editor_upload", controller.RichEditorUpload)

	//查询分类
	r.GET("/category", controller.SearchCategory)
	r.GET("/category/:id", controller.SearchCategoryName)

	//用户信息管理
	userRoutes := r.Group("/user")
	userRoutes.Use(middleware.AuthMiddleware())
	//验证用户信息
	userRoutes.GET("", controller.GetInfo)
	//获取用户简要信息
	userRoutes.GET("briefInfo/:id", controller.GetBriefInfo)
	//获取用户详细信息
	userRoutes.GET("detailedInfo/:id", controller.GetDetailInfo)
	//修改头像
	userRoutes.PUT("avatar/:id", controller.ModifyAvatar)
	//修改用户名
	userRoutes.PUT("name/:id", controller.ModavatarifyName)

	//收藏
	colRoutes := r.Group("/collects")
	colRoutes.Use(middleware.AuthMiddleware())
	colRoutes.GET(":id", controller.FindCollects)
	colRoutes.PUT("new/:id", controller.NewCollect)
	colRoutes.DELETE(":index", controller.UnCollect)

	//关注
	folRoutes := r.Group("/following")
	folRoutes.Use(middleware.AuthMiddleware())
	folRoutes.GET(":id", controller.FindFollowing)
	folRoutes.PUT("new/:id", controller.NewFollowing)
	folRoutes.DELETE(":index", controller.UnFollowing)

	//用户文章的增删改查
	articleRoutes := r.Group("/article")
	articleController := controller.NewArticleController()
	//增
	articleRoutes.POST("", middleware.AuthMiddleware(), articleController.Create)
	//删
	articleRoutes.DELETE(":id", middleware.AuthMiddleware(), articleController.Delete)
	//改
	articleRoutes.PUT(":id", middleware.AuthMiddleware(), articleController.Update)
	//查
	articleRoutes.GET(":id", articleController.Show)
	articleRoutes.POST("list", articleController.List)

	return r
}
