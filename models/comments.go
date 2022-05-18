package models

import "gorm.io/gorm"

// Comments 用户评论文章数据库模型
type Comments struct {
	gorm.Model
	ID               int `gorm:"primaryKey;autoIncrement"`
	ArticleID        int
	CommentingUserID int
	Content          string `gorm:"size:200"`
}
