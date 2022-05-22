package controller

import (
	"gin-blog/common"
	"gin-blog/form/user"
	"gin-blog/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//RegisterGET 注册页面GET请求
func RegisterGET(ctx *gin.Context) {
	ctx.HTML(200, "user/register.html", gin.H{
		"message": "welcome!",
	})

}

//RegisterPOST 注册页面POST请求
func RegisterPOST(ctx *gin.Context) {
	db := common.GetDB()
	var register user.RegisterForm
	//获取登录参数
	err := ctx.ShouldBind(&register)
	if err != nil {
		ctx.HTML(422, "user/register.html", gin.H{
			"message": err.Error(),
		})
		return
	}
	//获取用户名,密码,手机号
	username := register.Username
	password := register.Password
	phone := register.PhoneNumber
	//判断数据长度 5<username<=20 5<password<=20 phone = 11
	if len(username) <= 5 && len(username) > 20 {
		ctx.HTML(422, "user/register.html", gin.H{
			"message": "用户名长度范围为:5-20!",
		})
		return
	}
	if len(password) <= 5 && len(password) > 20 {
		ctx.HTML(422, "user/register.html", gin.H{
			"message": "密码长度范围为:5-20!",
		})
		return
	}
	if len(phone) != 11 {
		ctx.HTML(422, "user/register.html", gin.H{
			"message": "手机号码长度必须为:11!",
		})
		return
	}
	//用户名查重
	var user models.Users
	db.Where("username = ?", username).First(&user)
	if user.ID != 0 {
		ctx.HTML(422, "user/register.html", gin.H{"message": "用户已经存在!不可重复注册!"})
		return
	}
	//密码使用sha256加密
	timeNow := time.Now().UnixNano()
	//将timeNow int64转换为string
	salt := strconv.FormatInt(timeNow, 10)
	hashPwd := common.GetHashPassword(password, salt)
	//写入数据库
	AddUserInfo := models.Users{Username: username, HashPassword: hashPwd, Salt: salt, Phone: phone, Active: true}
	db.Create(&AddUserInfo)
	//返回HTML
	ctx.HTML(200, "user/register.html", gin.H{
		"message": "注册账号: " + AddUserInfo.Username + " 成功!",
	})

}
