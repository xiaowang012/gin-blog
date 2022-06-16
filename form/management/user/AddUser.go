package user

// AddUser 后台管理添加用户表单结构体
type AddUser struct {
	PhoneNumber string `form:"phone" binding:"required"`
	Username    string `form:"username" binding:"required"`
	Password    string `form:"password" binding:"required"`
}
