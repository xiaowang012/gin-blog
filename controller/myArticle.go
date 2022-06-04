package controller

import (
	"fmt"
	"gin-blog/common"
	"gin-blog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//MyArticle 我的文章页面
func MyArticle(ctx *gin.Context) {
	db := common.GetDB()
	//获取当前用户
	session := sessions.Default(ctx)
	currentUserInfo := session.Get("currentUser").(UserInfo)
	//查询以当前用户为作者的文章，限制返回10条
	var articles []models.Articles
	db.Where("author = ?", currentUserInfo.UserName).Order("created_at desc").Find(&articles).Limit(10)
	//返回数据到HTML
	ctx.HTML(200, "my/my.html", gin.H{
		"currentUser": currentUserInfo.UserName,
		"msg":         fmt.Sprintf("文章列表: %d 篇文章", len(articles)),
		"style":       "alert alert-success alert-dismissable",
		"articles":    articles,
		"currentPage": 1,
	})
}

// MyArticlePage 我的文章界面分页
func MyArticlePage(ctx *gin.Context) {
	db := common.GetDB()
	//获取当前用户
	session := sessions.Default(ctx)
	currentUserInfo := session.Get("currentUser").(UserInfo)
	//获取页码参数
	pageNumber := ctx.Query("pageNumber")
	//将pageNumber转换为int
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		return
	}
	//定义接收mysql article，messages数据的slice
	var articles []models.Articles
	//查询数据库
	db.Where("author = ?", currentUserInfo.UserName).Order("created_at desc").Limit(pageNumberInt * 10).Offset((pageNumberInt - 1) * 10).Find(&articles)
	ctx.HTML(200, "my/my.html", gin.H{
		"currentUser": currentUserInfo.UserName,
		"msg":         fmt.Sprintf("文章列表: %d 篇文章", len(articles)),
		"style":       "alert alert-success alert-dismissable",
		"articles":    articles,
		"currentPage": pageNumberInt,
	})

}

// MyArticleDelete 我的文章界面删除文章
func MyArticleDelete(ctx *gin.Context) {
	//我的文章删除文章
	db := common.GetDB()
	session := sessions.Default(ctx)
	currentUserInfo := session.Get("currentUser")
	currentUser := currentUserInfo.(UserInfo).UserName
	//获取文章id
	id := ctx.Query("id")
	//将id转换为int
	deleteID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		return
	}
	//查询是否有文章
	var article models.Articles
	db.First(&article, deleteID)
	if article.ID == 0 {
		fmt.Println("要删除的文章不存在!无法删除")
		ctx.Redirect(http.StatusTemporaryRedirect, "/my")
		return
	}
	//执行删除文章操作
	//验证作者和当前用户是否一致，否者无权限删除
	if article.Author != currentUser {
		fmt.Println("删除失败! 无法删除其他用户的文章!")
		ctx.Redirect(http.StatusTemporaryRedirect, "/my")
		return
	}

	db.Delete(&article, article.ID)
	ctx.Redirect(http.StatusTemporaryRedirect, "/my")

}

//我的文章界面搜索文章

//我的文章界面搜索文章分页

//我的文章界面修改文章
