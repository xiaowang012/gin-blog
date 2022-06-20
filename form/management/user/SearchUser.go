package user

// SearchUser 后台管理用户管理查询用户(根据用户名查询)
type SearchUser struct {
	UserName string `form:"username" binding:"required"`
}
