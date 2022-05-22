package controller

import (
	"gin-blog/common"
	"gin-blog/models"
	"github.com/gin-gonic/gin"
)

func IndexGET(ctx *gin.Context) {
	db := common.GetDB()
	//获取当前登录用户
	curUser := "wangli"
	//查询文章前五条数据,使用切片接收多条数据
	var articles []models.Articles
	db.Limit(5).Offset(0).Find(&articles)
	//查询结果为 slice
	//for _, article := range articles {
	//	oldTime := article.ReleaseDate
	//	TimeFormat := oldTime.Format("2006-01-02 ")
	//
	//}
	//查留言板前5条数据
	var messages []models.MessageBoard
	db.Limit(5).Offset(0).Find(&messages)
	ctx.HTML(200, "user/index.html", gin.H{
		"message":     "ok",
		"currentUser": curUser,
		"msg":         curUser + " ,欢迎来到Blog!",
		//"style":   "alert alert-dismissible alert-danger",
		"style":    "alert alert-success alert-dismissable",
		"articles": articles,
		"messages": messages,
	})

}
