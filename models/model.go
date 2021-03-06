package models

import (
	"gorm.io/gorm"
)

// Users 用户注册数据库模型
type Users struct {
	gorm.Model
	ID           int    `gorm:"primaryKey;autoIncrement"`
	Username     string `gorm:"size:100;unique"`
	Nickname     string `gorm:"size:100"`
	HashPassword string `gorm:"size:100"`
	Salt         string `gorm:"size:100"`
	Email        string `gorm:"size:100"`
	PicturePath  string `gorm:"size:100"`
	RegisterTime string `gorm:"size:100"`
	Birthday     string `gorm:"size:100"`
	Age          int    `gorm:"size:100"`
	Phone        string `gorm:"size:100"`
	Active       bool
}

// Articles 博客文章数据库模型
type Articles struct {
	gorm.Model
	ID                  int    `gorm:"primaryKey;autoIncrement"`
	ReleaseDate         string `gorm:"size:300"`
	Author              string `gorm:"size:300"`
	BlogTitle           string `gorm:"size:300"`
	BlogContentOverview string `gorm:"size:300"`
	BlogContent         string `gorm:"text"`
	Likes               int
	Comments            int
	NumberOfViews       int
	IfAnonymous         bool
	BlogPicturePath     string `gorm:"size:600"`
	Tag                 string `gorm:"size:100"`
}

// Comments 用户评论文章数据库模型
type Comments struct {
	gorm.Model
	ID             int `gorm:"primaryKey;autoIncrement"`
	ArticleID      int
	CommentingUser string `gorm:"size:300"`
	Content        string `gorm:"size:600"`
	IfAnonymous    bool
}

// Group 用户组数据库模型
type Group struct {
	gorm.Model
	ID        int    `gorm:"primaryKey;autoIncrement"`
	GroupName string `gorm:"size:100"`
}

// Permissions 路由权限数据库模型
type Permissions struct {
	gorm.Model
	ID          int    `gorm:"primaryKey;autoIncrement"`
	Url         string `gorm:"size:100"`
	GroupName   string `gorm:"size:100"`
	Description string `gorm:"size:100"`
}

// UserGroup 用户-角色数据库模型
type UserGroup struct {
	gorm.Model
	ID      int `gorm:"primaryKey;autoIncrement"`
	UserID  int
	GroupID int
}

// MessageBoard 用户首页留言版数据库模型
type MessageBoard struct {
	gorm.Model
	ID          int    `gorm:"primaryKey;autoIncrement"`
	PostUser    string `gorm:"size:100"`
	Content     string `gorm:"size:600"`
	IfAnonymous bool
}

// ArticleTags 文章分类标签 数据库模型
type ArticleTags struct {
	gorm.Model
	ID  int    `gorm:"primaryKey;autoIncrement"`
	Tag string `gorm:"size:100;unique"`
}
