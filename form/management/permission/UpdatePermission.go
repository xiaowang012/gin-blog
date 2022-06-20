package permission

// UpdatePermission 后台管理权限管理修改权限
type UpdatePermission struct {
	ID          string `form:"update_id" binding:"required"`
	Url         string `form:"update_url" binding:"required"`
	GroupName   string `form:"update_group_name" binding:"required"`
	Description string `form:"update_description" binding:"required"`
}
