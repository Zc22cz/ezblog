package controller

import (
	"blog_server/common"
	"blog_server/model"
	"blog_server/response"
	"blog_server/vo"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

type ArticleController struct {
	DB *gorm.DB
}

type IArticleController interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Show(c *gin.Context)
	List(c *gin.Context)
}

// 创建文章
func (a ArticleController) Create(c *gin.Context) {
	var articleRequest vo.CreatArticleRequest

	if err := c.ShouldBindJSON(&articleRequest); err != nil {
		response.Fail(c, nil, "数据错误")
		return
	}
	user, _ := c.Get("user")
	article := model.Article{
		UserId:     user.(model.User).ID,
		CategoryId: articleRequest.CategoryId,
		Title:      articleRequest.Title,
		Content:    articleRequest.Content,
		HeadImage:  articleRequest.HeadImage,
	}

	if err := a.DB.Create(&article).Error; err != nil {
		response.Fail(c, nil, "创建失败")
		return
	}
	response.Success(c, gin.H{"id": article.ID}, "创建成功")
}

// 更新文章
func (a ArticleController) Update(c *gin.Context) {
	var articleRequest vo.CreatArticleRequest
	err := c.ShouldBindJSON(&articleRequest)
	if err != nil {
		response.Fail(c, nil, "数据错误")
		return
	}
	//查id
	articleId := c.Params.ByName("id")
	var article model.Article
	//查文章
	err = a.DB.Where("id=?", articleId).First(&article).Error
	if err != nil {
		response.Fail(c, nil, "文章不存在")
		return
	}
	//获取用户
	user, _ := c.Get("user")
	userId := user.(model.User).ID
	if userId != article.UserId {
		response.Fail(c, nil, "用户不符合")
		return
	}
	//更新文章
	err = a.DB.Model(&article).Update(articleRequest).Error
	if err != nil {
		response.Fail(c, nil, "更新失败")
		return
	}

	response.Success(c, nil, "更新成功")
}

// 删除文章
func (a ArticleController) Delete(c *gin.Context) {
	//获取id
	articleId := c.Params.ByName("id")
	//查文章
	var article model.Article
	err := a.DB.Where("id=?", articleId).First(&article).Error
	if err != nil {
		response.Fail(c, nil, "文章不存在")
		return
	}
	//验证用户
	user, _ := c.Get("user")
	userId := user.(model.User).ID
	if userId != article.UserId {
		response.Fail(c, nil, "用户不符合")
		return
	}
	//删除
	err = a.DB.Delete(&article).Error
	if err != nil {
		response.Fail(c, nil, "删除失败")
		return
	}

	response.Success(c, nil, "删除成功")
}

// 展示单篇文章
func (a ArticleController) Show(c *gin.Context) {
	articleId := c.Params.ByName("id")
	var article model.Article
	err := a.DB.Where("id=?", articleId).First(&article).Error
	if err != nil {
		response.Fail(c, nil, "文章不存在")
		return
	}

	response.Success(c, gin.H{"article": article}, "查找成功")
}

// 以分页的形式展示多篇文章
func (a ArticleController) List(c *gin.Context) {
	//获取参数
	keyword := c.DefaultQuery("keyword", "")
	categoryId := c.DefaultQuery("categoryId", "0")
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "3"))
	var query []string
	var args []string
	//判断参数是否存在
	if keyword != "" {
		query = append(query, "(title LIKE ? OR content LIKE ?)")
		args = append(args, "%"+keyword+"%")
		args = append(args, "%"+keyword+"%")
	}
	if categoryId != "0" {
		query = append(query, "category_id = ?")
		args = append(args, categoryId)
	}
	//拼接字符串
	var queryStr string
	if len(query) > 0 {
		queryStr = strings.Join(query, "AND")
	}

	var article []model.ArticleInfo
	var count int
	offsetVal := (pageNum - 1) * pageSize

	//查询文章
	switch len(args) {
	case 0:
		a.DB.Table("articles").
			Select("id, category_id, title, LEFT(content, 80), head_image, created_at").
			Order("created_at desc").Offset(offsetVal).Limit(pageSize).Find(&article)
		a.DB.Model(model.Article{}).Count(&count)
	case 1:
		a.DB.Table("articles").
			Select("id, category_id, title, LEFT(content, 80), head_image, created_at").
			Where(queryStr, args[0]).Order("created_at desc").Offset(offsetVal).
			Limit(pageSize).Find(&article)
		a.DB.Model(model.Article{}).Where(queryStr, args[0]).Count(&count)
	case 2:
		a.DB.Table("articles").
			Select("id, category_id, title, LEFT(content, 80), head_image, created_at").
			Where(queryStr, args[0], args[1]).Order("created_at desc").Offset(offsetVal).
			Limit(pageSize).Find(&article)
		a.DB.Model(model.Article{}).Where(queryStr, args[0], args[1]).Count(&count)
	case 3:
		a.DB.Table("articles").
			Select("id, category_id, title, LEFT(content, 80), head_image, created_at").
			Where(queryStr, args[0], args[1], args[2]).Order("created_at desc").Offset(offsetVal).
			Limit(pageSize).Find(&article)
		a.DB.Model(model.Article{}).Where(queryStr, args[0], args[1], args[2]).Count(&count)
	}
	response.Success(c, gin.H{"article": article, "count": count}, "查找成功")
}

func NewArticleController() IArticleController {
	db := common.GetDB()
	db.AutoMigrate(model.Article{})
	return ArticleController{DB: db}
}
