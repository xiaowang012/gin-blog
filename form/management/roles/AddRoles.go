package roles

// AddRoles 后台管理角色管理添加角色
type AddRoles struct {
	RoleName string `form:"role_name" binding:"required"`
}
