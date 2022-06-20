package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-blog/common"
	"gin-blog/form/management/permission"
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
	//获取当前登录用户
	session := sessions.Default(ctx)
	currentUserInfo := session.Get("currentUser")
	currentUser := currentUserInfo.(UserInfo).UserName
	//获取redis中用户管理首页的用户数据切片
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
	//获取db连接
	db := common.GetDB()
	var searchUserInfo user.SearchUser
	//获取登录参数
	err = ctx.ShouldBind(&searchUserInfo)
	//表单出错
	if err != nil {
		ctx.HTML(422, "management/user.html", gin.H{
			"currentUser": currentUser,
			"msg":         "错误: " + err.Error(),
			"style":       "alert alert-dismissible alert-danger",
			"users":       users,
			"currentPage": 1,
		})
		return
	}
	//获取查询参数
	UserName := searchUserInfo.UserName
	//查表
	var searchUsers []models.Users
	db.Where(fmt.Sprintf(" username like %q ", "%"+UserName+"%")).Limit(10).Find(&searchUsers)
	ctx.HTML(200, "management/user_search.html", gin.H{
		"currentUser": currentUser,
		"msg":         fmt.Sprintf("查询成功! 查询到%d条数据", len(searchUsers)),
		"style":       "alert alert-success alert-dismissable",
		"users":       searchUsers,
		"currentPage": 1,
		"kw":          UserName,
	})
}

//UserManagementSearchUsersPage 后台管理用户管理界面搜索用户分页
func UserManagementSearchUsersPage(ctx *gin.Context) {
	//获取db连接
	db := common.GetDB()
	//获取页码参数以及查询参数
	searchInfo := ctx.Query("search")
	pageNumber := ctx.Query("page")
	//将pageNumber转换为int
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		return
	}
	//查询用户信息
	var searchUsers []models.Users
	db.Where(fmt.Sprintf(" username like %q ", "%"+searchInfo+"%")).Limit(pageNumberInt * 10).Offset((pageNumberInt - 1) * 10).Find(&searchUsers)
	//返回数据到HTML
	ctx.HTML(200, "management/user_search.html", gin.H{
		"msg":         fmt.Sprintf("查询成功! 查询到%d条数据", len(searchUsers)),
		"style":       "alert alert-success alert-dismissable",
		"users":       searchUsers,
		"currentPage": pageNumberInt,
		"kw":          searchInfo,
	})

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
	//获取redis连接
	rdb := common.GetRedis()
	//获取db连接
	db := common.GetDB()
	//获取当前登录用户
	session := sessions.Default(ctx)
	currentUserInfo := session.Get("currentUser").(UserInfo)
	//查询权限表的数据，限制返回10条
	var permissionsInfo []models.Permissions
	db.Limit(10).Order("created_at desc").Find(&permissionsInfo)
	//判断redis中是否存在key management_permission_page_1
	_, err := rdb.Get(context.Background(), "management_permission_page_1").Bytes()
	if err != nil {
		byteData, err := json.Marshal(permissionsInfo)
		if err != nil {
			fmt.Println("permission转换数据错误: " + err.Error())
		} else {
			err = rdb.Set(context.Background(), "management_permission_page_1", byteData, time.Minute*10).Err()
			if err != nil {
				println("SET KEY: management_permission_page_1错误:" + err.Error())
			}
		}
	}
	//返回数据到HTML
	ctx.HTML(200, "management/permission.html", gin.H{
		"currentUser": currentUserInfo.UserName,
		"permissions": permissionsInfo,
		"currentPage": 1,
	})

}

//PermissionManagementPage 后台管理权限管理页面分页
func PermissionManagementPage(ctx *gin.Context) {
	db := common.GetDB()
	session := sessions.Default(ctx)
	//获取当前登录用户
	userinfo := session.Get("currentUser")
	userinfoNew := userinfo.(UserInfo)
	//获取页码参数
	pageNumber := ctx.Query("page")
	//将pageNumber转换为int
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		return
	}
	//查询 permissions
	var permissionsInfo []models.Permissions
	//查询数据库
	db.Limit(pageNumberInt * 10).Offset((pageNumberInt - 1) * 10).Find(&permissionsInfo)
	//返回数据到HTML
	ctx.HTML(200, "management/permission.html", gin.H{
		"currentUser": userinfoNew.UserName,
		"permissions": permissionsInfo,
		"currentPage": pageNumberInt,
	})

}

//PermissionManagementAddPermission 后台管理权限管理添加权限
func PermissionManagementAddPermission(ctx *gin.Context) {
	//获取redis 中的权限数据
	var permissionsInfo []models.Permissions
	rdb := common.GetRedis()
	val, err := rdb.Get(context.Background(), "management_permission_page_1").Bytes()
	if err != nil {
		fmt.Println("读取management_permission_page_1失败!")
	} else {
		err = json.Unmarshal(val, &permissionsInfo)
		if err != nil {
			fmt.Println("解析management_permission_page_1错误:" + err.Error())
		}
	}
	//获取当前登录用户
	session := sessions.Default(ctx)
	currentUserInfo := session.Get("currentUser").(UserInfo)
	//获取db
	db := common.GetDB()
	var addPermission permission.AddPermission
	//绑定添加权限表单
	err = ctx.ShouldBind(&addPermission)
	if err != nil {
		ctx.HTML(422, "management/permission.html", gin.H{
			"style":       "alert alert-dismissible alert-danger",
			"msg":         "表单错误: " + err.Error(),
			"permissions": permissionsInfo,
			"currentUser": currentUserInfo.UserName,
			"currentPage": 1,
		})
		return
	}
	//获取URL,GroupName,Description
	Url := addPermission.Url
	GroupName := addPermission.GroupName
	Description := addPermission.Description

	//判断数据长度
	if len(Url) > 100 {
		ctx.HTML(422, "management/permission.html", gin.H{
			"style":       "alert alert-dismissible alert-danger",
			"msg":         "Url长度最大为:100!",
			"permissions": permissionsInfo,
			"currentUser": currentUserInfo.UserName,
			"currentPage": 1,
		})
		return
	}
	if len(GroupName) > 100 {
		ctx.HTML(422, "management/permission.html", gin.H{
			"style":       "alert alert-dismissible alert-danger",
			"msg":         "用户组最大长度为:100!",
			"permissions": permissionsInfo,
			"currentUser": currentUserInfo.UserName,
			"currentPage": 1,
		})
		return
	}
	if len(Description) > 100 {
		ctx.HTML(422, "management/permission.html", gin.H{
			"style":       "alert alert-dismissible alert-danger",
			"msg":         "描述信息最大长度为:100!",
			"permissions": permissionsInfo,
			"currentUser": currentUserInfo.UserName,
			"currentPage": 1,
		})
		return
	}
	//写入数据库
	AddPermissionInfo := models.Permissions{
		Url:         Url,
		GroupName:   GroupName,
		Description: Description,
	}
	db.Create(&AddPermissionInfo)
	//返回HTML
	ctx.HTML(200, "management/permission.html", gin.H{
		"style": "alert alert-success alert-dismissable",
		"msg": fmt.Sprintf("添加权限: %s %s %s  成功!",
			addPermission.Url,
			addPermission.GroupName,
			addPermission.Description,
		),
		"currentPage": 1,
		"currentUser": currentUserInfo.UserName,
		"permissions": permissionsInfo,
	})

}

//PermissionManagementUpdatePermission 后台管理权限管理修改权限
func PermissionManagementUpdatePermission(ctx *gin.Context) {
	//获取redis 中的权限数据
	var permissionsInfo []models.Permissions
	rdb := common.GetRedis()
	val, err := rdb.Get(context.Background(), "management_permission_page_1").Bytes()
	if err != nil {
		fmt.Println("读取management_permission_page_1失败!")
	} else {
		err = json.Unmarshal(val, &permissionsInfo)
		if err != nil {
			fmt.Println("解析management_permission_page_1错误:" + err.Error())
		}
	}
	//获取当前登录用户
	session := sessions.Default(ctx)
	currentUserInfo := session.Get("currentUser").(UserInfo)
	//获取db
	db := common.GetDB()
	var updatePermission permission.UpdatePermission
	//绑定添加权限表单
	err = ctx.ShouldBind(&updatePermission)
	if err != nil {
		ctx.HTML(422, "management/permission.html", gin.H{
			"style":       "alert alert-dismissible alert-danger",
			"msg":         "表单错误: " + err.Error(),
			"permissions": permissionsInfo,
			"currentUser": currentUserInfo.UserName,
			"currentPage": 1,
		})
		return
	}
	//获取ID,URL,GroupName,Description
	ID := updatePermission.ID
	Url := updatePermission.Url
	GroupName := updatePermission.GroupName
	Description := updatePermission.Description

	//将ID转换为int
	IDint, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		return
	}

	//查询permission
	var permissionInfo models.Permissions
	db.First(&permissionInfo, IDint)
	if permissionInfo.ID == 0 {
		fmt.Println("permission不存在! 无法修改!")
		return
	}

	//判断数据长度
	if len(Url) > 100 {
		ctx.HTML(422, "management/permission.html", gin.H{
			"style":       "alert alert-dismissible alert-danger",
			"msg":         "Url长度最大为:100!",
			"permissions": permissionsInfo,
			"currentUser": currentUserInfo.UserName,
			"currentPage": 1,
		})
		return
	}
	if len(GroupName) > 100 {
		ctx.HTML(422, "management/permission.html", gin.H{
			"style":       "alert alert-dismissible alert-danger",
			"msg":         "用户组最大长度为:100!",
			"permissions": permissionsInfo,
			"currentUser": currentUserInfo.UserName,
			"currentPage": 1,
		})
		return
	}
	if len(Description) > 100 {
		ctx.HTML(422, "management/permission.html", gin.H{
			"style":       "alert alert-dismissible alert-danger",
			"msg":         "描述信息最大长度为:100!",
			"permissions": permissionsInfo,
			"currentUser": currentUserInfo.UserName,
			"currentPage": 1,
		})
		return
	}
	//修改字段
	db.Model(&permissionInfo).Updates(models.Permissions{
		Url:         Url,
		GroupName:   GroupName,
		Description: Description,
	})
	//返回HTML
	ctx.HTML(200, "management/permission.html", gin.H{
		"style":       "alert alert-success alert-dismissable",
		"msg":         "修改数据成功!",
		"currentPage": 1,
		"currentUser": currentUserInfo.UserName,
		"permissions": permissionsInfo,
	})

}

//PermissionManagementDeletePermission 后台管理权限管理删除权限
func PermissionManagementDeletePermission(ctx *gin.Context) {
	//获取db连接
	db := common.GetDB()
	//获取用户ID参数
	PermissionId := ctx.Query("id")
	//将id 由string转换为int
	PermissionIdInt, err := strconv.Atoi(PermissionId)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		ctx.Redirect(http.StatusMovedPermanently, "/management/permission")
		return
	}
	//根据userId 查询用户
	var permissionInfo models.Permissions
	db.First(&permissionInfo, PermissionIdInt)
	//判断用户是否存在
	if permissionInfo.ID == 0 {
		fmt.Println("用户不存在! 无法删除此账号!")
		ctx.Redirect(http.StatusMovedPermanently, "/management/permission")
		return
	}
	//删除用户
	db.Delete(&permissionInfo)
	ctx.Redirect(http.StatusMovedPermanently, "/management/permission")

}
