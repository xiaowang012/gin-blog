package models

import (
	"gorm.io/gorm"
)

// Group 用户组数据库模型
type Group struct {
	gorm.Model
	ID        int    `gorm:"primaryKey;autoIncrement"`
	GroupName string `gorm:"size:100"`
}
