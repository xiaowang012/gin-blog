package article

//AddCommentsForm 给文章添加评论信息
type AddCommentsForm struct {
	UserName  string `form:"username" binding:"required"`
	ArticleID string `form:"articleID" binding:"required"`
	Content   string `form:"content" binding:"required"`
	Anonymous string `form:"anonymous"`
}
