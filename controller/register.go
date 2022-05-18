package controller

import (
	"gin-blog/form/user"
	"gin-blog/models"

	"github.com/gin-gonic/gin"
)

func register(ctx *gin.Context) {
	var register user.RegisterForm
	//获取登录参数
	err := ctx.ShouldBind(&register)
	if err != nil {
		ctx.JSON(422, gin.H{"message": err.Error()})
		return
	} else {
		username := register.Username
		password := register.Password
		phone := register.PhoneNumber
		//判断数据长度 5<username<=20 5<password<=20 phone = 11
		if len(username) <= 5 && len(username) > 20 {
			ctx.JSON(422, gin.H{"message": "用户名长度范围为:5-20!"})
			return
		}
		if len(password) <= 5 && len(password) > 20 {
			ctx.JSON(422, gin.H{"message": "密码长度范围为:5-20!"})
			return
		}
		if len(phone) != 11 {
			ctx.JSON(422, gin.H{"message": "手机号码长度必须为:11!"})
			return
		}
		//密码使用sha256加密
		//hashPwd := GetHashPassword(password)
		//fmt.Println(hashPwd)
		//写入数据库
		active := true
		user := models.Users{Username: username, HashPassword: password, Phone: phone, Active: active}
		//db.Create(&user)
		//返回创建数据的ID
		//fmt.Println(user.ID)
		//返回json
		ctx.JSON(200, gin.H{
			"message":  "register success!",
			"userid":   user.ID,
			"username": user.Username,
		})
	}
}
