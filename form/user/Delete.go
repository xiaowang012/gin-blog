package user

//删除用户表单
type DeleteUserForm struct {
	Username string `form:"username",binding:"required"`
}
