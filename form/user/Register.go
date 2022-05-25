package user

// RegisterForm 注册表单结构体
type RegisterForm struct {
	PhoneNumber string `form:"phone" binding:"required"`
	Username    string `form:"username" binding:"required"`
	Password    string `form:"password" binding:"required"`
}
