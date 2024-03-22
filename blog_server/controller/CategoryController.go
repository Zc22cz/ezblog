package controller

import (
	"blog_server/common"
	"blog_server/model"
	"blog_server/response"
	"github.com/gin-gonic/gin"
)

// 查询分类
func SearchCategory(c *gin.Context) {
	db := common.GetDB()
	var categories []model.Category
	err := db.Find(&categories).Error
	if err != nil {
		response.Fail(c, nil, "查找失败")
		return
	}
	response.Success(c, gin.H{"categories": categories}, "查找成功")
}

// 查询分类id
func SearchCategoryName(c *gin.Context) {
	db := common.GetDB()
	var category model.Category
	categoryId := c.Param("id")
	err := db.Where("id=?", categoryId).First(&category).Error
	if err != nil {
		response.Fail(c, nil, "该分类不存在")
		return
	}
	response.Success(c, gin.H{"categoryName": category.CategoryName}, "查找成功")
}
