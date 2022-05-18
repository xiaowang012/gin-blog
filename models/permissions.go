package models

import (
	"gorm.io/gorm"
)

// Permissions 路由权限数据库模型
type Permissions struct {
	gorm.Model
	ID          int    `gorm:"primaryKey;autoIncrement"`
	Url         string `gorm:"size:100"`
	GroupName   string `gorm:"size:100"`
	Description string `gorm:"size:100"`
}
