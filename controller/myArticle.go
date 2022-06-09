package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-blog/common"
	"gin-blog/form/article"
	"gin-blog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

//MyArticle 我的文章页面
func MyArticle(ctx *gin.Context) {
	db := common.GetDB()
	//获取当前用户
	session := sessions.Default(ctx)
	currentUserInfo := session.Get("currentUser").(UserInfo)
	//查询以当前用户为作者的文章，限制返回10条
	var articles []models.Articles
	db.Where("author = ?", currentUserInfo.UserName).Limit(10).Order("created_at desc").Find(&articles)
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

// MyArticleSearch 我的文章界面搜索文章
func MyArticleSearch(ctx *gin.Context) {
	//获取当前用户
	session := sessions.Default(ctx)
	currentUserInfo := session.Get("currentUser")
	currentUser := currentUserInfo.(UserInfo).UserName
	//获取db连接
	db := common.GetDB()
	var searchInfo article.SearchArticlesForm
	//获取登录参数
	err := ctx.ShouldBind(&searchInfo)
	//表单出错
	if err != nil {
		var articles []models.Articles
		db.Where("author = ?", currentUser).Order("created_at desc").Find(&articles).Limit(10)
		ctx.HTML(200, "my/my.html", gin.H{
			"currentUser": currentUser,
			"msg":         "错误: 查询信息不能为空!",
			"style":       "alert alert-dismissible alert-danger",
			"articles":    articles,
			"currentPage": 1,
		})
		return
	}
	//获取查询参数
	name := searchInfo.ArticleName
	//查表
	var articles []models.Articles
	db.Where(fmt.Sprintf(" blog_title like %q ", "%"+name+"%")).Limit(10).Find(&articles)
	ctx.HTML(200, "my/my_search.html", gin.H{
		"currentUser": currentUser,
		"msg":         fmt.Sprintf("查询成功! 查询到%d条数据", len(articles)),
		"style":       "alert alert-success alert-dismissable",
		"articles":    articles,
		"currentPage": 1,
		"kw":          name,
	})
}

// MyArticleSearchPage 我的文章界面搜索文章分页
func MyArticleSearchPage(ctx *gin.Context) {
	db := common.GetDB()
	//获取页码参数以及查询参数
	name := ctx.Query("search")
	pageNumber := ctx.Query("pageNumber")
	//将pageNumber转换为int
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		return
	}
	//定义接收mysql article，messages数据的slice
	var articles []models.Articles
	db.Where(fmt.Sprintf(" blog_title like %q ", "%"+name+"%")).Limit(pageNumberInt * 10).Offset((pageNumberInt - 1) * 10).Find(&articles)
	//返回数据到HTML
	ctx.HTML(200, "my/my_search.html", gin.H{
		"msg":         fmt.Sprintf("查询成功! 查询到%d条数据", len(articles)),
		"style":       "alert alert-success alert-dismissable",
		"articles":    articles,
		"currentPage": pageNumberInt,
		"kw":          name,
	})
}

// MyArticleUpdateArticlePage 我的文章界面修改文章页面
func MyArticleUpdateArticlePage(ctx *gin.Context) {
	//获取db
	db := common.GetDB()
	//获取当前用户
	session := sessions.Default(ctx)
	CurrentUserInfo := session.Get("currentUser")
	currentUser := CurrentUserInfo.(UserInfo).UserName
	//获取文章id
	articleID := ctx.Query("id")
	//根据ID查询文章
	var articleData models.Articles
	db.First(&articleData, articleID)
	//获取tags
	var tags []models.ArticleTags
	rdb := common.GetRedis()
	val, err := rdb.Get(context.Background(), "tags").Bytes()
	if err != nil {
		fmt.Println("获取tagsInfo失败：" + err.Error())
		//从数据库中获取tags
		fmt.Println("从数据库中获取tags-------------")
		db.Find(&tags)
		byteData, err := json.Marshal(tags)
		if err != nil {
			fmt.Println("tags 转换数据错误: " + err.Error())
		}
		err = rdb.Set(context.Background(), "tags", byteData, time.Hour*2).Err()
		if err != nil {
			println("SET tags错误!:" + err.Error())
		}
	} else {
		err = json.Unmarshal(val, &tags)
		if err != nil {
			fmt.Println("解析tags失败-----")
			return
		}

	}
	var tagsInfo []string
	for _, v := range tags {
		tagsInfo = append(tagsInfo, v.Tag)
	}
	//拼接html
	var option string
	var optionHtml []string
	for _, v := range tagsInfo {
		option = fmt.Sprintf("<option value=\"%s\">%s</option>", v, v)
		//fmt.Println(option)
		optionHtml = append(optionHtml, option)
	}
	ctx.HTML(200, "article/edit.html", gin.H{
		"ID":          articleData.ID,
		"currentUser": currentUser,
		"Author":      articleData.Author,
		"title":       articleData.BlogTitle,
		"Overview":    articleData.BlogContentOverview,
		"op":          articleData.Tag,
		"Content":     articleData.BlogContent,
		"options":     optionHtml,
	})

}

// MyArticleUpdateArticle 我的文章界面修改文章POST请求
func MyArticleUpdateArticle(ctx *gin.Context) {
	//获取db
	db := common.GetDB()
	//获取当前用户
	session := sessions.Default(ctx)
	currentUser := session.Get("currentUser").(UserInfo)
	//从数据库或者redis中获取tags
	//渲染tags 标签
	var tags []models.ArticleTags
	rdb := common.GetRedis()
	val, err := rdb.Get(context.Background(), "tags").Bytes()
	if err != nil {
		fmt.Println("获取tagsInfo失败：" + err.Error())
		//从数据库中获取tags
		fmt.Println("从数据库中获取tags-------------")
		db.Find(&tags)
		byteData, err := json.Marshal(tags)
		if err != nil {
			fmt.Println("tags 转换数据错误: " + err.Error())
		}
		err = rdb.Set(context.Background(), "tags", byteData, time.Hour*2).Err()
		if err != nil {
			println("SET tags错误!:" + err.Error())
		}
	} else {
		err = json.Unmarshal(val, &tags)
		if err != nil {
			fmt.Println("解析tags失败-----")
			return
		}

	}
	var tagsInfo []string
	for _, v := range tags {
		tagsInfo = append(tagsInfo, v.Tag)
	}
	//拼接html
	var option string
	var optionHtml []string
	for _, v := range tagsInfo {
		option = fmt.Sprintf("<option value=\"%s\">%s</option>", v, v)
		optionHtml = append(optionHtml, option)
	}
	var articleInfo article.UpdateArticleForm
	//获取POST请求中的文章数据
	err = ctx.ShouldBind(&articleInfo)
	if err != nil {
		ctx.HTML(422, "article/edit.html", gin.H{
			"msg":     "错误! " + err.Error(),
			"style":   "alert alert-dismissible alert-danger",
			"Author":  currentUser.UserName,
			"options": optionHtml,
		})
		return
	}
	//获取ID,作者，标题,内容，等参数
	ID := articleInfo.ID
	Author := articleInfo.Author
	Title := articleInfo.BlogTitle
	Overview := articleInfo.BlogContentOverview
	Content := articleInfo.Content
	Tag := articleInfo.Tag

	//将ID由string转换为int
	IdInt, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		return
	}
	//根据ID查询文章是否存在
	var articleUpdate models.Articles
	db.First(&articleUpdate, IdInt)
	//判断是否查询到文章
	if articleUpdate.ID == 0 {
		fmt.Println("未查询到文章! 无法修改文章信息!")
		return
	}
	//判断数据长度是否符合要求
	if len(Author) > 300 {
		ctx.HTML(422, "article/edit.html", gin.H{
			"msg":     "错误! 用户名长度范围最大为:100!",
			"style":   "alert alert-dismissible alert-danger",
			"Author":  currentUser.UserName,
			"options": optionHtml,
		})
		return
	}
	if len(Title) > 300 {
		ctx.HTML(422, "article/edit.html", gin.H{
			"msg":     "错误! 标题长度范围最大为:100!",
			"style":   "alert alert-dismissible alert-danger",
			"Author":  currentUser.UserName,
			"options": optionHtml,
		})
		return
	}
	if len(Overview) > 300 {
		ctx.HTML(422, "article/edit.html", gin.H{
			"msg":     "错误! 文章概述长度范围最大为:100!",
			"style":   "alert alert-dismissible alert-danger",
			"Author":  currentUser.UserName,
			"options": optionHtml,
		})
		return
	}
	if len(Content) > 16383 {
		ctx.HTML(422, "article/edit.html", gin.H{
			"msg":     "错误! 内容长度范围最大为:16383!",
			"style":   "alert alert-dismissible alert-danger",
			"Author":  currentUser.UserName,
			"options": optionHtml,
		})
		return
	}

	//判断作者是否为currentUser一致
	if currentUser.UserName != Author {
		ctx.HTML(422, "article/edit.html", gin.H{
			"msg":      "failed:发布失败!不允许修改用户名提交!",
			"style":    "alert alert-dismissible alert-danger",
			"Author":   currentUser.UserName,
			"title":    Title,
			"Overview": Overview,
			"Content":  Content,
			"options":  optionHtml,
		})
		return

	}
	//更新字段
	db.Model(&articleUpdate).Updates(models.Articles{
		Author:              Author,
		BlogTitle:           Title,
		BlogContentOverview: Overview,
		BlogContent:         Content,
		Tag:                 Tag,
	})

	//返回HTML信息
	ctx.HTML(200, "article/edit.html", gin.H{
		"msg":     "success:修改文章成功!可以到主页查看或者搜索查看修改后的文章!",
		"style":   "alert alert-success alert-dismissable",
		"Author":  currentUser.UserName,
		"options": optionHtml,
	})
}
