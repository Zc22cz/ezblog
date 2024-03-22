package model

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Article struct {
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primary_key;"`
	UserId     uint      `json:"user_id" gorm:"not null"`
	CategoryId uint      `json:"category_id" gorm:"not null"`
	Title      string    `json:"title" gorm:"type:varchar(50);not null"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	HeadImage  string    `json:"head_image"`
	CreatedAt  Time      `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt  Time      `json:"updated_at" gorm:"type:timestamp;default:'2023-08-08 17:50:05'"`
}

type ArticleInfo struct {
	ID         string `json:"id"`
	CategoryId uint   `json:"category_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	HeadImage  string `json:"head_image"`
	CreatedAt  Time   `json:"created_at" gorm:"type:timestamp;default:'2023-08-08 17:50:05'"`
}

// BeforeCreate 在创建文章之前将id赋值
// 通常情况下，我们使用 *gorm.DB 来执行数据库操作，比如查询、更新、删除等
// 通过 s *gorm.Scope，我们可以设置模型的一些属性，比如设置主键值等
func (a *Article) BeforeCreate(s *gorm.Scope) error {
	return s.SetColumn("ID", uuid.NewV4())
}
