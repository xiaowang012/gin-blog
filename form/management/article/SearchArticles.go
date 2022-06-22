package article

// SearchArticles 后台管理文章管理查询文章
type SearchArticles struct {
	BlogTitle string `form:"blog_title" binding:"required"`
}
