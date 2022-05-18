package models

import (
	"gorm.io/gorm"
	"time"
)

// Articles 博客文章数据库模型
type Articles struct {
	gorm.Model
	ID              int `gorm:"primaryKey;autoIncrement"`
	ReleaseDate     time.Time
	Author          string `gorm:"size:100"`
	BlogTitle       string `gorm:"size:100"`
	BlogContent     string `gorm:"size:500"`
	Likes           int
	comments        int
	BlogPicturePath string `gorm:"size:200"`
}
