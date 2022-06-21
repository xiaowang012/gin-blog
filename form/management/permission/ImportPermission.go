package permission

// ImportPermission 后台管理权限管理导入权限
type ImportPermission struct {
	Url string `form:"url" binding:"required"`
}
