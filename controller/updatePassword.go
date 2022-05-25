package controller

import (
	"gin-blog/common"
	"gin-blog/form/user"
	"gin-blog/models"
	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

//ChangePasswordGET 修改密码GET请求
func ChangePasswordGET(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userinfo := session.Get("currentUser")
	var currentUser string
	if userinfo == nil {
		currentUser = ""
	} else {
		userInfoStruct := userinfo.(UserInfo)
		currentUser = userInfoStruct.UserName
	}
	ctx.HTML(200, "user/pwd.html", gin.H{
		"currentUser": currentUser,
		"message":     "welcome!",
	})

}

//ChangePasswordPOST 修改密码POST请求
func ChangePasswordPOST(ctx *gin.Context) {
	db := common.GetDB()
	var changePassword user.ChangePasswordForm
	//获取登录参数
	err := ctx.ShouldBind(&changePassword)
	if err != nil {
		ctx.HTML(422, "user/pwd.html", gin.H{
			"message": err.Error(),
		})
		return
	}
	//获取用户名,密码,手机号
	username := changePassword.Username
	oldPassword := changePassword.OldPassword
	newPassword := changePassword.NewPassword
	//判断数据长度 5<username<=20 5<password<=20 phone = 11
	if len(username) <= 5 && len(username) > 20 {
		ctx.HTML(422, "user/pwd.html", gin.H{
			"message": "用户名长度范围为:5-20!",
		})
		return
	}
	if len(oldPassword) <= 5 && len(oldPassword) > 20 {
		ctx.HTML(422, "user/pwd.html", gin.H{
			"message": "原密码长度范围为:5-20!",
		})
		return
	}
	if len(newPassword) <= 5 && len(newPassword) > 20 {
		ctx.HTML(422, "user/pwd.html", gin.H{
			"message": "新密码长度范围为:5-20!",
		})
		return
	}
	//查询要修改的用户是否存在
	var user models.Users
	db.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		ctx.HTML(422, "user/pwd.html", gin.H{"message": "用户不存在!无法修改密码!"})
		return
	}
	//验证原密码是否正确
	salt := user.Salt
	pwd := user.HashPassword
	//使用数据库中获取的salt与获取的 oldPassword sha256加密
	hashPwd := common.GetHashPassword(oldPassword, salt)
	//对比两个密码是否一致
	if pwd != hashPwd {
		//返回HTML
		ctx.HTML(403, "user/pwd.html", gin.H{
			"message": " 原密码错误!无法修改!",
		})
		return
	}
	//将新密码写入数据库
	newHashPwd := common.GetHashPassword(newPassword, salt)
	//db.Model(&models.Users{}).Where("username = ?", username).Update("HashPassword", hashPwd)
	db.Model(&user).Update("HashPassword", newHashPwd)
	//返回HTML
	ctx.HTML(200, "user/pwd.html", gin.H{
		"message": "修改账号: " + user.Username + "的密码成功!",
	})

}
