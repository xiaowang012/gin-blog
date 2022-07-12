package roles

// SearchRoles 后台管理角色管理查询角色
type SearchRoles struct {
	RoleName string `form:"role_name" binding:"required"`
}
