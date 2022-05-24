package controller

import (
	"fmt"
	"gin-blog/common"
	"gin-blog/form/article"
	"gin-blog/models"
	"github.com/gin-gonic/gin"
)

func SearchArticles(ctx *gin.Context) {
	db := common.GetDB()
	//session := sessions.Default(ctx)
	////获取当前登录用户
	//userinfo := session.Get("currentUser")
	//if userinfo == nil {
	//	ctx.Redirect(http.StatusMovedPermanently, "/login")
	//	return
	//}
	//userinfoNew := userinfo.(UserInfo)
	////判断UserInfo数据是否为空
	//if userinfoNew.UserName == "" || userinfoNew.ExpirationTime == "" {
	//	ctx.Redirect(http.StatusMovedPermanently, "/login")
	//	return
	//}
	////判断session id中的时间是否过期
	//ExpirationTime := userinfoNew.ExpirationTime
	//CurrentTime := time.Now().Format("2006-01-02 15:04:05")
	////先把时间字符串格式化成相同的时间类型
	//t1, err := time.Parse("2006-01-02 15:04:05", ExpirationTime)
	//t2, err := time.Parse("2006-01-02 15:04:05", CurrentTime)
	//if err == nil && t1.Before(t2) {
	//	//session失效，清空session，UserInfo 重定向到login页面
	//	session.Delete("currentUser")
	//	session.Save()
	//	ctx.Redirect(http.StatusMovedPermanently, "/login")
	//	return
	//}
	//curUser := userinfoNew.UserName
	var searchInfo article.SearchArticlesForm
	//获取登录参数
	err1 := ctx.ShouldBind(&searchInfo)
	//表单出错
	if err1 != nil {
		var articles []models.Articles
		db.Limit(5).Offset(0).Find(&articles)
		//查留言板前5条数据
		var messages []models.MessageBoard
		db.Limit(5).Offset(0).Find(&messages)
		ctx.HTML(200, "user/index.html", gin.H{
			"currentUser": "",
			"msg":         "错误: 查询信息不能为空!",
			"style":       "alert alert-dismissible alert-danger",
			"articles":    articles,
			"messages":    messages,
		})
		return
	}
	//获取用户名,密码,手机号
	SearchArticleName := searchInfo.ArticleName
	//查表
	var articles []models.Articles
	db.Where(fmt.Sprintf(" blog_title like %q ", "%"+SearchArticleName+"%")).Limit(5).Find(&articles)
	ctx.HTML(200, "user/index.html", gin.H{
		"currentUser": curUser,
		"msg":         fmt.Sprintf("查询成功! 查询到%d条数据", len(articles)),
		"style":       "alert alert-success alert-dismissable",
		"articles":    articles,
	})
}
