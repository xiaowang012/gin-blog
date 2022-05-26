package middleware

//
//import (
//	"github.com/gin-contrib/sessions"
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"time"
//)
//
////LoginRequired 检查session是否有值，是否为空，是否过期，返回当前用户
//func LoginRequired(ctx *gin.Context) {
//	session := sessions.Default(ctx)
//	//获取当前登录用户(结构体)
//	userinfo := session.Get("currentUser")
//	if userinfo == nil {
//		//未获取到session，重定向到login
//		ctx.Redirect(http.StatusMovedPermanently, "/login")
//		return
//	}
//	userinfoNew := userinfo.(UserInfo)
//	//判断UserInfo数据是否为空
//	if userinfoNew.UserName == "" || userinfoNew.ExpirationTime == "" {
//		//session结构体中的数据为空
//		ctx.Redirect(http.StatusMovedPermanently, "/login")
//		return
//	}
//	//判断session id中的时间是否过期
//	ExpirationTime := userinfoNew.ExpirationTime
//	CurrentTime := time.Now().Format("2006-01-02 15:04:05")
//	//先把时间字符串格式化成相同的时间类型
//	t1, err := time.Parse("2006-01-02 15:04:05", ExpirationTime)
//	t2, err := time.Parse("2006-01-02 15:04:05", CurrentTime)
//	if err == nil && t1.Before(t2) {
//		//session失效，清空session，UserInfo 重定向到login页面
//		session.Delete("currentUser")
//		session.Save()
//		return
//	}
//	curUser := userinfoNew.UserName
//	ctx.Set("currentUser", curUser)
//	return
//
//}
