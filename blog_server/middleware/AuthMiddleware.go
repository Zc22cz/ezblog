package middleware

import (
	"blog_server/common"
	"blog_server/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取authorization header
		tokenString := c.Request.Header.Get("Authorization")
		//token为空
		if tokenString == "" {
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			c.Abort()
			return
		}
		//非法token
		if tokenString == "" || len(tokenString) < 7 || !strings.HasPrefix(tokenString, "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			c.Abort()
			return
		}
		//提取token有效部分
		tokenString = tokenString[7:]
		//解析token
		token, claims, err := common.ParseToken(tokenString)
		//非法token
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限错误",
			})
			c.Abort()
			return
		}
		//获取claims中的userid
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.Where("id=?", userId).First(&user)
		//将用户信息存入上下文以便获取
		c.Set("user", user)
		c.Next()
	}
}
