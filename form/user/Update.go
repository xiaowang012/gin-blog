package user

//修改密码表单
type ChangePasswordForm struct {
	Username    string `form:"username",binding:"required"`
	OldPassword string `form:"oldpassword",binding:"required"`
	NewPassword string `form:"newpassword",binding:"required"`
}
