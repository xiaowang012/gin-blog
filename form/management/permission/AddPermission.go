package permission

// AddPermission 后台管理权限管理添加权限
type AddPermission struct {
	Url         string `form:"url" binding:"required"`
	GroupName   string `form:"group_name" binding:"required"`
	Description string `form:"description" binding:"required"`
}
