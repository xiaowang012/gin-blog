package user

// UpdateUser 后台管理修改用户信息表单结构体
type UpdateUser struct {
	ID       string `form:"update_id" binding:"required"`
	NickName string `form:"update_nickname" binding:"required"`
	Email    string `form:"update_email" binding:"required"`
	Birthday string `form:"update_birthday" binding:"required"`
	Age      string `form:"update_age" binding:"required"`
	Phone    string `form:"update_phone" binding:"required"`
}
