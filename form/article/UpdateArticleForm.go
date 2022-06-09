package article

//UpdateArticleForm 修改文章提交表单
type UpdateArticleForm struct {
	ID                  string `form:"ID" binding:"required"`
	Author              string `form:"Author" binding:"required"`
	BlogTitle           string `form:"BlogTitle" binding:"required"`
	BlogContentOverview string `form:"BlogContentOverview" binding:"required"`
	Content             string `form:"summernote" binding:"required"`
	Tag                 string `form:"tag" binding:"required"`
}
