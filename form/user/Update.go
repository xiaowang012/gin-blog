package user

//修改密码表单
type ChangePasswordForm struct {
	Username    string `form:"username" binding:"required"`
	OldPassword string `form:"password" binding:"required"`
	NewPassword string `form:"password1" binding:"required"`
}
