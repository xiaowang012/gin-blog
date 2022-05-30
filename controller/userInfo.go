package controller

import (
	"fmt"
	"gin-blog/common"
	"gin-blog/form/userinfo"
	"gin-blog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

//UserInfoPage 用户个人信息页面
func UserInfoPage(ctx *gin.Context) {
	db := common.GetDB()
	session := sessions.Default(ctx)
	//获取当前登录用户
	userinfo := session.Get("currentUser")
	if userinfo == nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	userinfoNew := userinfo.(UserInfo)
	//判断UserInfo数据是否为空
	if userinfoNew.UserName == "" || userinfoNew.ExpirationTime == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	//判断session id中的时间是否过期
	ExpirationTime := userinfoNew.ExpirationTime
	CurrentTime := time.Now().Format("2006-01-02 15:04:05")
	//先把时间字符串格式化成相同的时间类型
	t1, err1 := time.Parse("2006-01-02 15:04:05", ExpirationTime)
	t2, err2 := time.Parse("2006-01-02 15:04:05", CurrentTime)
	if err1 == nil && err2 == nil && t1.Before(t2) {
		//session失效，清空session，UserInfo 重定向到login页面
		session.Delete("currentUser")
		session.Save()
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	//根据当前用户名查询用户个人信息数据
	var user models.Users
	db.Where("username = ?", userinfoNew.UserName).First(&user)
	if user.ID == 0 {
		//用户不存在
		ctx.Redirect(http.StatusTemporaryRedirect, "/index")
		return
	}
	//fmt.Println(user.PicturePath)
	ctx.HTML(200, "user/userinfo.html", gin.H{
		"msg":         "",
		"ID":          user.ID,
		"Nickname":    user.Nickname,
		"Username":    user.Username,
		"Age":         user.Age,
		"Birthday":    user.Birthday,
		"Phone":       user.Phone,
		"Email":       user.Email,
		"PicturePath": user.PicturePath,
	})

}

//UserInfoUpdate 用户界面编辑个人信息
func UserInfoUpdate(ctx *gin.Context) {
	db := common.GetDB()
	var updateUserInfo userinfo.ChangeUserInfoForm
	//获取提交修改的数据
	err := ctx.ShouldBind(&updateUserInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//获取参数
	ID := updateUserInfo.ID
	Username := updateUserInfo.Username
	Nickname := updateUserInfo.Nickname
	Email := updateUserInfo.Email
	Birthday := updateUserInfo.Birthday
	Age := updateUserInfo.Age
	Phone := updateUserInfo.Phone
	//获取头像文件
	var filePath string
	file, err := ctx.FormFile("picture")
	if err != nil {
		fmt.Println("接收文件错误:" + err.Error())
	} else {
		//保存到指定文件
		fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ".png"
		curDir, err := os.Getwd()
		if err != nil {
			fmt.Println("获取当前工作目录失败: " + err.Error())
			return
		}

		filePath = "/static/imgs/user/" + fileName
		err = ctx.SaveUploadedFile(file, curDir+filePath)
		if err != nil {
			fmt.Println("保存头像文件错误：" + err.Error())
			return
		}

	}

	//判断数据长度
	if len(Username) > 100 {
		fmt.Println("Username 长度最大为100!")
		return
	}
	if len(Nickname) > 100 {
		fmt.Println("Nickname 长度最大为100!")
		return
	}
	if len(Email) > 100 {
		fmt.Println("Email 长度最大为100!")
		return
	}
	if len(Birthday) > 100 {
		fmt.Println("Birthday 长度最大为100!")
		return
	}
	if len(Phone) > 100 {
		fmt.Println("Phone 长度最大为100!")
		return
	}

	//根据ID查询用户信息
	var user models.Users
	db.First(&user, ID)
	if user.ID == 0 {
		fmt.Println("未查询到用户,无法更新!")
		return
	}
	//写入数据库
	if filePath == "" {
		db.Model(&user).Updates(models.Users{Nickname: Nickname, Email: Email,
			Birthday: Birthday, Age: Age, Phone: Phone})
	} else {
		db.Model(&user).Updates(models.Users{Nickname: Nickname, Email: Email, PicturePath: filePath,
			Birthday: Birthday, Age: Age, Phone: Phone})
	}
	ctx.Redirect(http.StatusTemporaryRedirect, "/index/userinfo")
}
