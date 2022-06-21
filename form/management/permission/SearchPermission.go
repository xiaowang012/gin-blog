package permission

// SearchPermission 后台管理权限管理查询权限
type SearchPermission struct {
	Url string `form:"url" binding:"required"`
}
