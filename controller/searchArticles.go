package controller

import (
	"fmt"
	"gin-blog/common"
	"gin-blog/form/article"
	"gin-blog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//SearchArticles 首页查询文章
func SearchArticles(ctx *gin.Context) {
	//获取session中的用户
	session := sessions.Default(ctx)
	//获取当前登录用户
	userinfo := session.Get("currentUser")
	if userinfo == nil {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}
	userinfoNew := userinfo.(UserInfo)
	//判断UserInfo数据是否为空
	if userinfoNew.UserName == "" || userinfoNew.ExpirationTime == "" {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}
	//判断session id中的时间是否过期
	ExpirationTime := userinfoNew.ExpirationTime
	CurrentTime := time.Now().Format("2006-01-02 15:04:05")
	//先把时间字符串格式化成相同的时间类型
	t1, err1 := time.Parse("2006-01-02 15:04:05", ExpirationTime)
	t2, err2 := time.Parse("2006-01-02 15:04:05", CurrentTime)
	if err1 == nil && err2 == nil && t1.Before(t2) {
		//session失效，清空session，UserInfo 重定向到login页面
		session.Delete("currentUser")
		session.Save()
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}
	//获取db连接
	db := common.GetDB()
	var searchInfo article.SearchArticlesForm
	//获取登录参数
	err3 := ctx.ShouldBind(&searchInfo)
	//表单出错
	if err3 != nil {
		var articles []models.Articles
		db.Limit(5).Offset(0).Find(&articles)
		//查留言板前5条数据
		var messages []models.MessageBoard
		db.Limit(5).Offset(0).Order("created_at desc").Find(&messages)
		//根据messages切片中的IfAnonymous 字段是否为true,然后复制到新切片
		var messagesNew []models.MessageBoard
		for _, val := range messages {
			if val.IfAnonymous == true {
				val.PostUser = "****"
			}
			messagesNew = append(messagesNew, val)
		}
		ctx.HTML(200, "user/index.html", gin.H{
			"currentUser": userinfoNew.UserName,
			"msg":         "错误: 查询信息不能为空!",
			"style":       "alert alert-dismissible alert-danger",
			"articles":    articles,
			"messages":    messagesNew,
			"currentPage": 1,
		})
		return
	}
	//获取用户名,密码,手机号
	SearchArticleName := searchInfo.ArticleName
	//查表
	var articles []models.Articles
	var messages []models.MessageBoard
	db.Limit(5).Offset(0).Order("created_at desc").Find(&messages)
	db.Where(fmt.Sprintf(" blog_title like %q ", "%"+SearchArticleName+"%")).Limit(5).Find(&articles)
	//根据messages切片中的IfAnonymous 字段是否为true,然后复制到新切片
	var messagesNew []models.MessageBoard
	for _, val := range messages {
		if val.IfAnonymous == true {
			val.PostUser = "****"
		}
		messagesNew = append(messagesNew, val)
	}
	ctx.HTML(200, "user/index.html", gin.H{
		"currentUser": userinfoNew.UserName,
		"msg":         fmt.Sprintf("查询成功! 查询到%d条数据", len(articles)),
		"style":       "alert alert-success alert-dismissable",
		"articles":    articles,
		"messages":    messagesNew,
		"currentPage": 1,
	})
}

////查询文章分页
//func SearchArticlesPage(ctx *gin.Context) {
//	//获取session中的用户
//	session := sessions.Default(ctx)
//	//获取当前登录用户
//	userinfo := session.Get("currentUser")
//	if userinfo == nil {
//		ctx.Redirect(http.StatusMovedPermanently, "/login")
//		return
//	}
//	userinfoNew := userinfo.(UserInfo)
//	//判断UserInfo数据是否为空
//	if userinfoNew.UserName == "" || userinfoNew.ExpirationTime == "" {
//		ctx.Redirect(http.StatusMovedPermanently, "/login")
//		return
//	}
//	//判断session id中的时间是否过期
//	ExpirationTime := userinfoNew.ExpirationTime
//	CurrentTime := time.Now().Format("2006-01-02 15:04:05")
//	//先把时间字符串格式化成相同的时间类型
//	t1, err1 := time.Parse("2006-01-02 15:04:05", ExpirationTime)
//	t2, err2 := time.Parse("2006-01-02 15:04:05", CurrentTime)
//	if err1 == nil && err2 == nil && t1.Before(t2) {
//		//session失效，清空session，UserInfo 重定向到login页面
//		session.Delete("currentUser")
//		session.Save()
//		ctx.Redirect(http.StatusMovedPermanently, "/login")
//		return
//	}
//	//获取db连接
//	db := common.GetDB()
//	var searchInfo article.SearchArticlesForm
//	//获取登录参数
//	err3 := ctx.ShouldBind(&searchInfo)
//	//表单出错
//	if err3 != nil {
//		var articles []models.Articles
//		db.Limit(5).Offset(0).Find(&articles)
//		//查留言板前5条数据
//		var messages []models.MessageBoard
//		db.Limit(5).Offset(0).Find(&messages)
//		ctx.HTML(200, "user/index.html", gin.H{
//			"currentUser": "",
//			"msg":         "错误: 查询信息不能为空!",
//			"style":       "alert alert-dismissible alert-danger",
//			"articles":    articles,
//			"messages":    messages,
//		})
//		return
//	}
//	//获取用户名,密码,手机号
//	SearchArticleName := searchInfo.ArticleName
//	//查表
//	var articles []models.Articles
//	db.Where(fmt.Sprintf(" blog_title like %q ", "%"+SearchArticleName+"%")).Limit(5).Find(&articles)
//	ctx.HTML(200, "user/index.html", gin.H{
//		"currentUser": userinfoNew.UserName,
//		"msg":         fmt.Sprintf("查询成功! 查询到%d条数据", len(articles)),
//		"style":       "alert alert-success alert-dismissable",
//		"articles":    articles,
//	})
//}
