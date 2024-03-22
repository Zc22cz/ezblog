package controller

import (
	"blog_server/common"
	"blog_server/model"
	"blog_server/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

// 注册功能
func Register(c *gin.Context) {
	db := common.GetDB()

	//获取参数
	var requestUser model.User
	c.Bind(&requestUser)
	userName := requestUser.UserName
	password := requestUser.Password
	phoneNumber := requestUser.PhoneNumber

	//验证数据是否已经存在
	var user model.User
	db.Where("phone_number=?", phoneNumber).First(&user)
	if user.ID != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 422,
			"msg":  "用户已存在",
		})
	}

	//密码加密
	newPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	//创建用户
	newUser := model.User{
		UserName:    userName,
		PhoneNumber: phoneNumber,
		Password:    string(newPassword),
		Avatar:      "/images/default_avater.png",
		Collects:    model.Array{},
		Following:   model.Array{},
		Fans:        0,
	}
	db.Create(&newUser)

	//返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
	})
}

// 登录功能
func Login(c *gin.Context) {
	db := common.GetDB()
	//获取参数
	var requestUser model.User
	c.Bind(&requestUser)
	phoneNumber := requestUser.PhoneNumber
	password := requestUser.Password
	//校验用户是否存在
	var user model.User
	db.Where("phone_number=?", phoneNumber).First(&user)
	if user.ID == 0 {
		c.JSON(200, gin.H{
			"code": 422,
			"msg":  "用户不存在",
		})
		return
	}
	//校验密码
	t := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if t != nil {
		c.JSON(200, gin.H{
			"code": 422,
			"msg":  "密码错误",
		})
		return
	}
	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "系统错误",
		})
		return
	}
	//返回结果
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{"token": token},
		"msg":  "登录成功",
	})
}

// 登录后获取信息
func GetInfo(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{"id": user.(model.User).ID, "avatar": user.(model.User).Avatar},
		"msg":  "登录信息获取成功",
	})
}

// 获取简要信息
func GetBriefInfo(c *gin.Context) {
	db := common.GetDB()
	userId := c.Params.ByName("id")
	user, _ := c.Get("user")
	var curUser model.User
	if userId == strconv.Itoa(int(user.(model.User).ID)) {
		curUser = user.(model.User)
	} else {
		db.Where("id=?", userId).First(&curUser)
		if curUser.ID == 0 {
			response.Fail(c, nil, "用户不存在")
			return
		}
	}
	response.Success(c, gin.H{
		"id":      curUser.ID,
		"name":    curUser.UserName,
		"avatar":  curUser.Avatar,
		"loginId": user.(model.User).ID,
	}, "查找成功")
}

// 获取详细信息
func GetDetailInfo(c *gin.Context) {
	db := common.GetDB()
	userId := c.Params.ByName("id")
	user, _ := c.Get("user")

	var curUser model.User
	if userId == strconv.Itoa(int(user.(model.User).ID)) {
		curUser = user.(model.User)
	} else {
		db.Where("id=?", userId).First(&curUser)
		if curUser.ID == 0 {
			response.Fail(c, nil, "用户不存在")
			return
		}
	}

	//返回用户详细信息
	var articles, collects []model.ArticleInfo
	var following []model.UserInfo
	var collist, follist []string
	collist = toStringArray(curUser.Collects)
	follist = toStringArray(curUser.Following)
	db.Table("articles").
		Select("id, category_id, title, LEFT(content, 80) AS content, head_image, created_at").
		Where("user_id=?", userId).Order("created_at desc").Find(&articles)
	db.Table("articles").
		Select("id, category_id, title, LEFT(content, 80) AS content, head_image, created_at").
		Where("id IN (?)", collist).Order("created_at desc").Find(&collects)
	db.Table("users").Select("id, avatar, user_name").Where("id IN (?)", follist).
		Find(&following)

	response.Success(c, gin.H{
		"id":        curUser.ID,
		"name":      curUser.UserName,
		"avatar":    curUser.Avatar,
		"loginId":   user.(model.User).ID,
		"articles":  articles,
		"collects":  collects,
		"following": following,
		"fans":      curUser.Fans,
	}, "查找成功")
}

func toStringArray(str []string) (a model.Array) {
	for i := 0; i < len(a); i++ {
		str = append(str, a[i])
	}
	return str
}

// 修改头像
func ModifyAvatar(c *gin.Context) {
	db := common.GetDB()
	user, _ := c.Get("user")
	var requestUser model.User
	c.Bind(&requestUser)
	avatar := requestUser.Avatar

	var curUser model.User
	db.Where("id=?", user.(model.User).ID).First(&curUser)

	//更新
	if err := db.Model(&curUser).Update("avatar", avatar).Error; err != nil {
		response.Fail(c, nil, "修改头像失败")
		return
	}

	response.Success(c, nil, "修改头像成功")
}

// 修改用户名
func ModavatarifyName(c *gin.Context) {
	db := common.GetDB()
	user, _ := c.Get("user")
	var requestUser model.User
	c.Bind(&requestUser)
	userName := requestUser.UserName

	var curUser model.User
	db.Where("id=?", user.(model.User).ID).First(&curUser)

	if err := db.Model(&curUser).Update("user_name", userName).Error; err != nil {
		response.Fail(c, nil, "修改用户名失败")
		return
	}
	response.Success(c, nil, "修改用户名成功")
}
