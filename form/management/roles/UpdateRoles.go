package roles

// UpdateRoles 后台管理角色管理修改角色名
type UpdateRoles struct {
	ID       string `form:"update_id" binding:"required"`
	RoleName string `form:"update_role_name" binding:"required"`
}
