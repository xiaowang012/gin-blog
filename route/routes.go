package route

import (
	"gin-blog/controller"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.Engine) *gin.Engine {
	//根据是否登录跳转到登录页或者主页
	//r.GET("/", controller.RegisterGET)
	//注册页面GET
	r.GET("/register", controller.RegisterGET)
	//注册界面POST
	r.POST("/register", controller.RegisterPOST)
	//登录页面GET
	r.GET("/login", controller.LoginGET)
	//登录页面POST
	r.POST("/login", controller.LoginPOST)

	//修改密码页面GET请求
	r.GET("/updatePassword", controller.ChangePasswordGET)
	//修改密码页面POST请求
	r.POST("/updatePassword", controller.ChangePasswordPOST)

	//用户主页GET
	r.GET("/index", controller.IndexGET)
	//用户主页POST

	return r
}
