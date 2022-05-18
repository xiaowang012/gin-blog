package models

import "gorm.io/gorm"

// UserGroup 用户-角色数据库模型
type UserGroup struct {
	gorm.Model
	ID      int `gorm:"primaryKey;autoIncrement"`
	UserID  int
	GroupID int
}

