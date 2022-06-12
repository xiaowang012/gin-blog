package controller

import (
	"gin-blog/common"
	"gin-blog/form/user"
	"gin-blog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//UserInfo session中的用户信息结构体
type UserInfo struct {
	UserName       string
	ExpirationTime string
}

//LoginGET 登录页面GET请求
func LoginGET(ctx *gin.Context) {
	ctx.HTML(200, "user/login.html", gin.H{
		"message": "welcome!",
	})

}

//LoginPOST 登录页面POST请求
func LoginPOST(ctx *gin.Context) {
	session := sessions.Default(ctx)
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
	//判断用户账号是否被禁用
	if user.Active == false {
		ctx.HTML(422, "user/login.html", gin.H{"message": " 用户账号已经被禁用! 请联系管理员解除! "})
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
	//设置session,过期时间设置为2小时
	curTime := time.Now()
	expirationTime := curTime.Add(time.Hour * 2).Format("2006-01-02 15:04:05")
	session.Set("currentUser", UserInfo{UserName: username, ExpirationTime: expirationTime})
	session.Save()
	ctx.Redirect(http.StatusMovedPermanently, "/index")

}
