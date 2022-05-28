package article

//WriteArticleForm 写文章提交表单结构体
type WriteArticleForm struct {
	Author              string `form:"Author" binding:"required"`
	BlogTitle           string `form:"BlogTitle" binding:"required"`
	BlogContentOverview string `form:"BlogContentOverview" binding:"required"`
	Anonymous           string `form:"anonymous"`
	Content             string `form:"summernote" binding:"required"`
}
