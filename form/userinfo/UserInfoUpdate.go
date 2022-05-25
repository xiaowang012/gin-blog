package userinfo

// ChangeUserInfoForm 修改用户个人信息
type ChangeUserInfoForm struct {
	ID       int    `form:"userid" binding:"required"`
	Username string `form:"username" binding:"required"`
	Nickname string `form:"nickname" binding:"required"`
	Email    string `form:"email" binding:"required"`
	Birthday string `form:"birthday" binding:"required"`
	Age      int    `form:"age" binding:"required"`
	Phone    string `form:"phone" binding:"required"`
}
