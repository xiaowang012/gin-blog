package models

import "gorm.io/gorm"

// UserLikes 用户点赞文章数据库模型
type Userlikes struct {
	gorm.Model
	ID         int `gorm:"primaryKey;autoIncrement"`
	ArticleID  int
	LikeUserID int
}
