package route

import (
	"gin-blog/controller"
	"github.com/gin-contrib/sessions"
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
	//退出登录 只能用HandleContext方式重定向所以卸载routes中
	r.GET("/logout", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		session.Delete("currentUser")
		session.Save()
		ctx.Request.URL.Path = "/login"
		r.HandleContext(ctx)
	})
	//修改密码
	r.GET("/updatePassword", controller.ChangePasswordGET)
	r.POST("/updatePassword", controller.ChangePasswordPOST)

	//Blog首页,首页展示，翻页，留言板
	r.GET("/index", controller.IndexGET)
	r.GET("/index/nextPage", controller.IndexGETNextPage)
	r.POST("/index/SendMessageBoard", controller.IndexMessageBoard)

	//用户个人信息页面
	r.GET("/index/userinfo", controller.UserInfoPage)
	//修改个人信息
	r.POST("/index/userinfo/update", controller.UserInfoUpdate)
	//搜索文章
	r.POST("/searchArticles", controller.SearchArticles)
	return r
}
