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
