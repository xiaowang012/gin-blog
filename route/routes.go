package route

import (
	"gin-blog/controller"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.Engine) *gin.Engine {
	//根据是否登录跳转到登录页或者主页
	//r.GET("/", controller.Host)

	//用户相关接口
	//注册
	r.GET("/register", controller.RegisterGET)
	r.POST("/register", controller.RegisterPOST)
	//登录
	r.GET("/login", controller.LoginGET)
	r.POST("/login", controller.LoginPOST)
	//退出登录 只能用临时重定向
	r.GET("/logout", controller.Logout)
	//修改密码
	r.GET("/updatePassword", controller.ChangePasswordGET)
	r.POST("/updatePassword", controller.ChangePasswordPOST)

	//Blog首页,首页展示，翻页，留言板
	r.GET("/index", controller.IndexGET)
	r.GET("/index/nextPage", controller.IndexGETNextPage)
	r.POST("/index/SendMessageBoard", controller.IndexMessageBoard)
	r.GET("/index/delete/messages", controller.IndexMessageDelete)

	//用户个人信息页面
	r.GET("/index/userinfo", controller.UserInfoPage)
	//修改个人信息
	r.POST("/index/userinfo/update", controller.UserInfoUpdate)
	//搜索文章
	r.POST("/searchArticles", controller.SearchArticles)
	//搜索文章分页
	r.GET("/searchArticles/page", controller.SearchArticlesPage)

	//文章处理
	r.GET("/article/details", controller.ArticleDetails)
	r.GET("/article/WriteArticle", controller.WriteArticlePage)
	r.POST("/article/WriteArticle", controller.WriteArticle)
	r.POST("/article/WriteArticle/picture/upload", controller.ReceivePicture)
	r.POST("/article/AddComments", controller.CommentingArticles)
	r.GET("/article/addLikes", controller.ArticleAddLikes)

	//文章列表页面相关路由
	r.GET("/article/list", controller.ArticleList)

	//我的文章
	r.GET("/my", controller.MyArticle)
	r.GET("/my/articles", controller.MyArticlePage)
	r.GET("/my/articles/Delete", controller.MyArticleDelete)
	r.POST("/my/articles/Search", controller.MyArticleSearch)
	r.GET("/my/articles/Search", controller.MyArticleSearchPage)
	r.GET("/my/articles/edit", controller.MyArticleUpdateArticlePage)
	r.POST("/my/articles/edit", controller.MyArticleUpdateArticle)

	//后台管理
	//用户管理路由
	r.GET("/management", controller.UserManagement)
	r.GET("/management/user/page", controller.UserManagementPage)
	r.GET("/management/user/disable", controller.UserManagementDisableUser)
	r.GET("management/user/enable", controller.UserManagementEnableUser)
	return r
}
