package controller

import (
	"fmt"
	"gin-blog/common"
	"gin-blog/form/user"
	"gin-blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

//LoginGET 登录页面GET请求
func LoginGET(ctx *gin.Context) {
	ctx.HTML(200, "user/login.html", gin.H{
		"message": "welcome!",
	})

}

//LoginPOST 登录页面POST请求
func LoginPOST(ctx *gin.Context) {
	fmt.Println("当前访问的URL: " + ctx.FullPath())
	db := common.GetDB()
	var login user.LoginForm
	//获取登录参数
	err := ctx.ShouldBind(&login)
	if err != nil {
		ctx.HTML(422, "user/login.html", gin.H{
			"message": err.Error(),
		})
		return
	}
	//获取用户名,密码
	username := login.Username
	password := login.Password
	//判断数据长度 5<username<=20 5<password<=20 phone = 11
	if len(username) <= 5 && len(username) > 20 {
		ctx.HTML(422, "user/login.html", gin.H{
			"message": "用户名长度范围为:5-20!",
		})
		return
	}
	if len(password) <= 5 && len(password) > 20 {
		ctx.HTML(422, "user/login.html", gin.H{
			"message": "密码长度范围为:5-20!",
		})
		return
	}
	//查询用户是否存在
	var user models.Users
	db.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		ctx.HTML(422, "user/login.html", gin.H{"message": " 用户不存在! "})
		return
	}
	//获取密码和salt
	salt := user.Salt
	pwd := user.HashPassword
	//使用数据库中获取的salt与新密码sha256加密
	hashPwd := common.GetHashPassword(password, salt)
	//对比两个密码是否一致
	if pwd != hashPwd {
		//返回HTML
		ctx.HTML(403, "user/login.html", gin.H{
			"message": " 账号或密码错误!",
		})
		return
	}

	//登录成功，重定向到index
	ctx.Redirect(http.StatusMovedPermanently, "/index")

}
