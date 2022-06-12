package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-blog/common"
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
	//将查询的数据存汝redis中
	//判断redis中是否存在key management_user10 (第一页的10条user数据)
	val, err := rdb.Get(context.Background(), "management_user10").Bytes()
	if err != nil {
		//fmt.Println("search mysql data --------------------")
		//查询数据库
		db.Limit(10).Order("created_at desc").Find(&users)
		byteData, err := json.Marshal(users)
		if err != nil {
			fmt.Println("management_user10转换数据错误: " + err.Error())
		}
		err = rdb.Set(context.Background(), "management_user10", byteData, time.Minute*10).Err()
		if err != nil {
			fmt.Println("SET management_user10错误!:" + err.Error())
		}
	} else {
		//fmt.Println("get redis data--------------")
		err = json.Unmarshal(val, &users)
		if err != nil {
			fmt.Println("解析users 数据失败!")
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
		return
	}
	//根据userId 查询用户
	var userInfo models.Users
	db.First(&userInfo, UserIdInt)
	//判断用户是否存在
	if userInfo.ID == 0 {
		fmt.Println("用户不存在! 无法禁用此账号!")
		return
	}
	//判断账号启用状态,如果已经处于禁用状态，则无需再次禁用
	if userInfo.Active == false {
		fmt.Println("用户已经处于禁用状态! 无法再次禁用此账号!")
		return
	}
	//将用户的Active 字段改为false
	db.Model(&userInfo).Update("active", false)
	ctx.Redirect(http.StatusMovedPermanently, "/management")

}

//UserManagementEnableUser 后台管理用户管理启用账号
func UserManagementEnableUser(ctx *gin.Context) {

}

//UserManagementUpdateUser 后台管理用户管理修改用户信息
func UserManagementUpdateUser(ctx *gin.Context) {

}

//UserManagementDeleteUser 后台管理用户管理删除用户
func UserManagementDeleteUser(ctx *gin.Context) {

}

//UserManagementAddUser 后台管理用户管理添加用户
func UserManagementAddUser(ctx *gin.Context) {

}

//UserManagementImportUsers 后台管理用户管理批量导入用户
func UserManagementImportUsers(ctx *gin.Context) {

}
