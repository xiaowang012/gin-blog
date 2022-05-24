package article

//SearchArticlesForm 查询文章表单结构体
type SearchArticlesForm struct {
	ArticleName string `form:"articleName" binding:"required"`
}
