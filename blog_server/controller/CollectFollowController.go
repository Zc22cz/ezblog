package controller

import (
	"blog_server/common"
	"blog_server/model"
	"blog_server/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 查询收藏
func FindCollects(c *gin.Context) {
	db := common.GetDB()
	user, _ := c.Get("user")
	id := c.Params.ByName("id")
	var curUser model.User
	db.Where("id=?", user.(model.User).ID).First(&curUser)

	for i := 0; i < len(curUser.Collects); i++ {
		if curUser.Collects[i] == id {
			response.Success(c, gin.H{"collected": true, "index": i}, "查询成功")
			return
		}
	}

	response.Success(c, gin.H{"collected": false}, "查询成功")
}

// 添加收藏
func NewCollect(c *gin.Context) {
	db := common.GetDB()
	user, _ := c.Get("user")
	id := c.Params.ByName("id")
	var curUser model.User
	db.Where("id=?", user.(model.User).ID).First(&curUser)
	var newCollects []string
	newCollects = append(curUser.Collects, id)

	if err := db.Model(&curUser).Update("collects", newCollects).Error; err != nil {
		response.Fail(c, nil, "添加失败")
		return
	}

	response.Success(c, nil, "添加成功")
}

// 删除收藏
func UnCollect(c *gin.Context) {
	db := common.GetDB()
	user, _ := c.Get("user")
	index, _ := strconv.Atoi(c.Params.ByName("index"))

	var curUser model.User
	db.Where("id=?", user.(model.User).ID).First(&curUser)
	var newCollects []string
	newCollects = append(curUser.Collects[:index], curUser.Collects[index+1:]...)

	if err := db.Model(&curUser).Update("collects", newCollects).Error; err != nil {
		response.Fail(c, nil, "删除失败")
		return
	}

	response.Success(c, nil, "删除成功")
}

// 查询关注
func FindFollowing(c *gin.Context) {
	db := common.GetDB()
	user, _ := c.Get("user")
	id := c.Params.ByName("id")
	var curUser model.User
	db.Where("id=?", user.(model.User).ID).First(&curUser)
	for i := 0; i < len(curUser.Following); i++ {
		if curUser.Following[i] == id {
			response.Success(c, gin.H{"followed": true, "index": i}, "查询成功")
			return
		}
	}
	response.Success(c, gin.H{"followed": false}, "查询成功")
}

// 添加关注
func NewFollowing(c *gin.Context) {
	db := common.GetDB()
	user, _ := c.Get("user")
	id := c.Params.ByName("id")
	var curUser model.User
	db.Where("id=?", user.(model.User).ID).First(&curUser)

	newFollowing := append(curUser.Following, id)
	if err := db.Model(&curUser).Update("following", newFollowing).Error; err != nil {
		response.Fail(c, nil, "添加失败")
		return
	}

	var followUer model.User
	db.Where("id=?", id).First(&followUer)
	if err := db.Model(&followUer).Update("fans", followUer.Fans+1).Error; err != nil {
		response.Fail(c, nil, "添加失败")
		return
	}

	response.Success(c, nil, "添加成功")
}

// 删除关注
func UnFollowing(c *gin.Context) {
	db := common.GetDB()
	user, _ := c.Get("user")
	index, _ := strconv.Atoi(c.Params.ByName("id"))
	var curUser model.User
	db.Where("id=?", user.(model.User).ID).First(&curUser)
	newFollowing := append(curUser.Following[:index], curUser.Following[index+1:]...)
	followId := curUser.Following[index]

	if err := db.Model(&curUser).Update("following", newFollowing).Error; err != nil {
		response.Fail(c, nil, "删除失败")
		return
	}

	var followUser model.User
	db.Where("id=?", followId).First(&followUser)
	if err := db.Model(&followUser).Update("fans", followUser.Fans-1).Error; err != nil {
		response.Fail(c, nil, "删除失败")
		return
	}

	response.Success(c, nil, "删除成功")
}
