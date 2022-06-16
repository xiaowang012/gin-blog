package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-blog/common"
	"gin-blog/form/management/user"
	"gin-blog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

//UserManagement 后台管理用户管理界面
func UserManagement(ctx *gin.Context) {
	//获取redis连接
	rdb := common.GetRedis()
	//获取db连接
	db := common.GetDB()
	//获取当前登录用户
	session := sessions.Default(ctx)
	currentUserInfo := session.Get("currentUser").(UserInfo)
	//查询用户表的数据，限制返回10条
	var users []models.Users
	db.Limit(10).Order("created_at desc").Find(&users)
	//判断redis中是否存在key management_user_page_1
	_, err := rdb.Get(context.Background(), "management_user_page_1").Bytes()
	if err != nil {
		byteData, err := json.Marshal(users)
		if err != nil {
			fmt.Println("users转换数据错误: " + err.Error())
		} else {
			err = rdb.Set(context.Background(), "management_user_page_1", byteData, time.Minute*10).Err()
			if err != nil {
				println("SET KEY: management_user_page_1错误:" + err.Error())
			}
		}
	}
	//返回数据到HTML
	ctx.HTML(200, "management/user.html", gin.H{
		"currentUser": currentUserInfo.UserName,
		"users":       users,
		"currentPage": 1,
	})
}

//UserManagementPage 后台管理用户管理界面分页
func UserManagementPage(ctx *gin.Context) {
	//获取db连接
	db := common.GetDB()
	//获取当前登录用户
	session := sessions.Default(ctx)
	currentUserInfo := session.Get("currentUser").(UserInfo)
	//获取页码参数
	pageNumber := ctx.Query("page")
	//将pageNumber 由string 转换为int
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		return

	}
	//查询用户表的数据，限制返回10条
	var users []models.Users
	db.Limit(pageNumberInt * 10).Offset((pageNumberInt - 1) * 10).Order("created_at desc").Find(&users)
	//返回数据到HTML
	ctx.HTML(200, "management/user.html", gin.H{
		"currentUser": currentUserInfo.UserName,
		"users":       users,
		"currentPage": pageNumberInt,
	})
}

//UserManagementSearchUsers 后台管理用户管理界面搜索用户
func UserManagementSearchUsers(ctx *gin.Context) {

}

//UserManagementSearchUsersPage 后台管理用户管理界面搜索用户分页
func UserManagementSearchUsersPage(ctx *gin.Context) {

}

//UserManagementDisableUser 后台管理用户管理禁用账号
func UserManagementDisableUser(ctx *gin.Context) {
	//获取db连接
	db := common.GetDB()
	//获取用户ID参数
	UserId := ctx.Query("id")
	//将id 由string转换为int
	UserIdInt, err := strconv.Atoi(UserId)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		ctx.Redirect(http.StatusMovedPermanently, "/management")
		return
	}
	//根据userId 查询用户
	var userInfo models.Users
	db.First(&userInfo, UserIdInt)
	//判断用户是否存在
	if userInfo.ID == 0 {
		fmt.Println("用户不存在! 无法禁用此账号!")
		ctx.Redirect(http.StatusMovedPermanently, "/management")
		return
	}
	//判断账号启用状态,如果已经处于禁用状态，则无需再次禁用
	if userInfo.Active == false {
		fmt.Println("用户已经处于禁用状态! 无法再次禁用此账号!")
		ctx.Redirect(http.StatusMovedPermanently, "/management")
		return
	}
	//将用户的Active 字段改为false
	db.Model(&userInfo).Update("active", false)
	ctx.Redirect(http.StatusMovedPermanently, "/management")

}

//UserManagementEnableUser 后台管理用户管理启用账号
func UserManagementEnableUser(ctx *gin.Context) {
	//获取db连接
	db := common.GetDB()
	//获取用户ID参数
	UserId := ctx.Query("id")
	//将id 由string转换为int
	UserIdInt, err := strconv.Atoi(UserId)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		ctx.Redirect(http.StatusMovedPermanently, "/management")
		return
	}
	//根据userId 查询用户
	var userInfo models.Users
	db.First(&userInfo, UserIdInt)
	//判断用户是否存在
	if userInfo.ID == 0 {
		fmt.Println("用户不存在! 无法启用此账号!")
		ctx.Redirect(http.StatusMovedPermanently, "/management")
		return
	}
	//判断账号启用状态,如果已经处于启用状态，则无需再次启用
	if userInfo.Active == true {
		fmt.Println("用户已经处于启用状态! 无法再次启用此账号!")
		ctx.Redirect(http.StatusMovedPermanently, "/management")
		return
	}
	//将用户的Active 字段改为false
	db.Model(&userInfo).Update("active", true)
	ctx.Redirect(http.StatusMovedPermanently, "/management")

}

//UserManagementUpdateUser 后台管理用户管理修改用户信息
func UserManagementUpdateUser(ctx *gin.Context) {
	//获取redis中存的后台管理用户管理第一页的数据
	var users []models.Users
	rdb := common.GetRedis()
	val, err := rdb.Get(context.Background(), "management_user_page_1").Bytes()
	if err != nil {
		fmt.Println("读取management_user_page_1失败!")
	} else {
		err = json.Unmarshal(val, &users)
		if err != nil {
			fmt.Println("解析messages错误:" + err.Error())
		}
	}
	//获取当前用户
	session := sessions.Default(ctx)
	currentUserInfo := session.Get("currentUser")
	currentUser := currentUserInfo.(UserInfo).UserName
	//获取db连接
	db := common.GetDB()
	var updateUser user.UpdateUser
	//绑定添加用户表单
	err = ctx.ShouldBind(&updateUser)
	if err != nil {
		ctx.HTML(422, "management/user.html", gin.H{
			"currentUser": currentUser,
			"msg":         "表单校验错误: " + err.Error(),
			"style":       "alert alert-dismissible alert-danger",
			"currentPage": 1,
			"users":       users,
		})
		return
	}
	//获取修改数据
	ID := updateUser.ID
	NickName := updateUser.NickName
	Email := updateUser.Email
	Birthday := updateUser.Birthday
	Age := updateUser.Age
	Phone := updateUser.Phone

	//将ID转换为int
	IdInt, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		return
	}
	//将Age转换为int
	AgeInt, err := strconv.Atoi(Age)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		return
	}

	//判断数据长度
	if len(NickName) > 100 {
		ctx.HTML(422, "management/user.html", gin.H{
			"currentUser": currentUser,
			"msg":         "用户名长度最大为: 100!",
			"style":       "alert alert-dismissible alert-danger",
			"currentPage": 1,
			"users":       users,
		})
		return
	}
	if len(Email) > 100 {
		ctx.HTML(422, "management/user.html", gin.H{
			"currentUser": currentUser,
			"msg":         "邮箱长度最大为: 100!",
			"style":       "alert alert-dismissible alert-danger",
			"currentPage": 1,
			"users":       users,
		})
		return
	}
	if len(Birthday) > 100 {
		ctx.HTML(422, "management/user.html", gin.H{
			"currentUser": currentUser,
			"msg":         "生日长度最大为: 100!",
			"style":       "alert alert-dismissible alert-danger",
			"currentPage": 1,
			"users":       users,
		})
		return
	}

	if len(Phone) != 11 {
		ctx.HTML(422, "management/user.html", gin.H{
			"currentUser": currentUser,
			"msg":         "手机号码长度必须为:11!",
			"style":       "alert alert-dismissible alert-danger",
			"currentPage": 1,
			"users":       users,
		})
		return
	}
	//查询用户是否存在
	var user models.Users
	db.First(&user, IdInt)
	if user.ID == 0 {
		ctx.HTML(422, "management/user.html", gin.H{
			"currentUser": currentUser,
			"msg":         "用户不存在! 无法修改用户信息!",
			"style":       "alert alert-dismissible alert-danger",
			"currentPage": 1,
			"users":       users,
		})
		return
	}
	//更新字段
	db.Model(&user).Updates(models.Users{Nickname: NickName, Email: Email,
		Birthday: Birthday, Age: AgeInt, Phone: Phone})
	//返回HTML
	ctx.HTML(200, "management/user.html", gin.H{
		"currentUser": currentUser,
		"msg":         "更新用户: " + user.Username + " 的信息成功!",
		"style":       "alert alert-success alert-dismissable",
		"currentPage": 1,
		"users":       users,
	})

}

//UserManagementDeleteUser 后台管理用户管理删除用户
func UserManagementDeleteUser(ctx *gin.Context) {
	//获取db连接
	db := common.GetDB()
	//获取用户ID参数
	UserId := ctx.Query("id")
	//将id 由string转换为int
	UserIdInt, err := strconv.Atoi(UserId)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		ctx.Redirect(http.StatusMovedPermanently, "/management")
		return
	}
	//根据userId 查询用户
	var userInfo models.Users
	db.First(&userInfo, UserIdInt)
	//判断用户是否存在
	if userInfo.ID == 0 {
		fmt.Println("用户不存在! 无法删除此账号!")
		ctx.Redirect(http.StatusMovedPermanently, "/management")
		return
	}
	//删除用户
	db.Delete(&userInfo)
	ctx.Redirect(http.StatusMovedPermanently, "/management")

}

//UserManagementAddUser 后台管理用户管理添加用户
func UserManagementAddUser(ctx *gin.Context) {
	//获取redis中存的后台管理用户管理第一页的数据
	var users []models.Users
	rdb := common.GetRedis()
	val, err := rdb.Get(context.Background(), "management_user_page_1").Bytes()
	if err != nil {
		fmt.Println("读取management_user_page_1失败!")
	} else {
		err = json.Unmarshal(val, &users)
		if err != nil {
			fmt.Println("解析messages错误:" + err.Error())
		}
	}
	//获取当前用户
	session := sessions.Default(ctx)
	currentUserInfo := session.Get("currentUser")
	currentUser := currentUserInfo.(UserInfo).UserName
	//获取db连接
	db := common.GetDB()
	var addUser user.AddUser
	//绑定添加用户表单
	err = ctx.ShouldBind(&addUser)
	if err != nil {
		ctx.HTML(422, "management/user.html", gin.H{
			"currentUser": currentUser,
			"msg":         "表单校验错误: " + err.Error(),
			"style":       "alert alert-dismissible alert-danger",
			"currentPage": 1,
			"users":       users,
		})
		return
	}
	//获取用户名,密码,手机号
	username := addUser.Username
	password := addUser.Password
	phone := addUser.PhoneNumber
	//判断数据长度 5<username<=20 5<password<=20 phone = 11
	if len(username) <= 5 && len(username) > 20 {
		ctx.HTML(422, "management/user.html", gin.H{
			"currentUser": currentUser,
			"msg":         "用户名长度范围为:5-20!",
			"style":       "alert alert-dismissible alert-danger",
			"currentPage": 1,
			"users":       users,
		})
		return
	}
	if len(password) <= 5 && len(password) > 20 {
		ctx.HTML(422, "management/user.html", gin.H{
			"currentUser": currentUser,
			"msg":         "密码长度范围为:5-20!",
			"style":       "alert alert-dismissible alert-danger",
			"currentPage": 1,
			"users":       users,
		})
		return
	}
	if len(phone) != 11 {
		ctx.HTML(422, "management/user.html", gin.H{
			"currentUser": currentUser,
			"msg":         "手机号码长度必须为:11!",
			"style":       "alert alert-dismissible alert-danger",
			"currentPage": 1,
			"users":       users,
		})
		return
	}
	//用户名查重
	var user models.Users
	db.Where("username = ?", username).First(&user)
	if user.ID != 0 {
		ctx.HTML(422, "management/user.html", gin.H{
			"currentUser": currentUser,
			"msg":         "用户已经存在!不可重复注册!",
			"style":       "alert alert-dismissible alert-danger",
			"currentPage": 1,
			"users":       users,
		})
		return
	}
	//密码使用sha256加密
	timeNow := time.Now().UnixNano()
	//将timeNow int64转换为string
	salt := strconv.FormatInt(timeNow, 10)
	hashPwd := common.GetHashPassword(password, salt)
	//写入数据库
	AddUserInfo := models.Users{
		Username:     username,
		HashPassword: hashPwd,
		Salt:         salt,
		Phone:        phone,
		Active:       true}
	db.Create(&AddUserInfo)
	//返回HTML
	ctx.HTML(200, "management/user.html", gin.H{
		"currentUser": currentUser,
		"msg":         "添加账号: " + AddUserInfo.Username + " 成功!",
		"style":       "alert alert-success alert-dismissable",
		"currentPage": 1,
		"users":       users,
	})

}

//PermissionManagement 后台管理权限管理页面
func PermissionManagement(ctx *gin.Context) {

}

//PermissionManagementPage 后台管理权限管理页面分页
func PermissionManagementPage(ctx *gin.Context) {

}
