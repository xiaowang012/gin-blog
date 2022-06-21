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
	"os"
	"strconv"
	"time"
)

// ArticlesInfoStruct redis缓存的文章编辑信息
type ArticlesInfoStruct struct {
	Title    string
	Overview string
	Content  string
}

//SearchArticles 首页查询文章
func SearchArticles(ctx *gin.Context) {
	//获取session中的用户
	session := sessions.Default(ctx)
	//获取当前登录用户
	userinfo := session.Get("currentUser")
	if userinfo == nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	userinfoNew := userinfo.(UserInfo)
	//判断UserInfo数据是否为空
	if userinfoNew.UserName == "" || userinfoNew.ExpirationTime == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
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
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
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
		db.Limit(5).Offset(0).Order("created_at desc").Find(&articles)
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
		ctx.HTML(200, "index/index.html", gin.H{
			"currentUser": userinfoNew.UserName,
			"msg":         "错误: 查询信息不能为空!",
			"style":       "alert alert-dismissible alert-danger",
			"articles":    articles,
			"messages":    messagesNew,
			"currentPage": 1,
		})
		return
	}
	//获取查询参数
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
	ctx.HTML(200, "index/index_search.html", gin.H{
		"currentUser": userinfoNew.UserName,
		"msg":         fmt.Sprintf("查询成功! 查询到%d条数据", len(articles)),
		"style":       "alert alert-success alert-dismissable",
		"articles":    articles,
		"messages":    messagesNew,
		"currentPage": 1,
		"kw":          SearchArticleName,
	})
}

//SearchArticlesPage 首页查询文章分页
func SearchArticlesPage(ctx *gin.Context) {
	db := common.GetDB()
	//获取页码参数以及查询参数
	searchInfo := ctx.Query("search")
	pageNumber := ctx.Query("pageNumber")
	//将pageNumber转换为int
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		return
	}
	//查询用户评论信息和文章数据
	//定义接收mysql article，messages数据的slice
	var articles []models.Articles
	var messages []models.MessageBoard
	db.Limit(5).Offset(0).Order("created_at desc").Find(&messages)
	db.Where(fmt.Sprintf(" blog_title like %q ", "%"+searchInfo+"%")).Limit(5).Offset((pageNumberInt - 1) * 5).Find(&articles)
	//根据messages切片中的IfAnonymous 字段是否为true,然后复制到新切片
	var messagesNew []models.MessageBoard
	for _, val := range messages {
		if val.IfAnonymous == true {
			val.PostUser = "****"
		}
		messagesNew = append(messagesNew, val)
	}
	//返回数据到HTML
	ctx.HTML(200, "index/index_search.html", gin.H{
		"msg":         fmt.Sprintf("查询成功! 查询到%d条数据", len(articles)),
		"style":       "alert alert-success alert-dismissable",
		"articles":    articles,
		"messages":    messagesNew,
		"currentPage": pageNumberInt,
		"kw":          searchInfo,
	})

}

// ArticleDetails 文章详情页面
func ArticleDetails(ctx *gin.Context) {
	session := sessions.Default(ctx)
	currentUser := session.Get("currentUser").(UserInfo)
	db := common.GetDB()
	articleID := ctx.Query("id")
	//将id 由string转为int
	articleIdInt, err := strconv.Atoi(articleID)
	if err != nil {
		fmt.Println("转换ID数据类型错误: " + err.Error())
		return
	}
	//查询文章信息
	var articleInfo models.Articles
	db.First(&articleInfo, articleIdInt)
	if articleInfo.ID == 0 {
		fmt.Println("未查询到文章信息,无法查看!")
		return
	}
	//根据文章ID查询文章的相关评论信息
	ArticleID := articleInfo.ID
	var CommentsInfo []models.Comments
	db.Where("article_id = ?", ArticleID).Limit(5).Order("created_at desc").Find(&CommentsInfo)
	//更新浏览量 +1
	numberOfViews := articleInfo.NumberOfViews
	numberOfViews++
	db.Model(&articleInfo).Updates(models.Articles{NumberOfViews: numberOfViews})

	//根据 CommentsInfo切片中的IfAnonymous 字段是否为true,然后复制到新切片
	var CommentsNew []models.Comments
	for _, val := range CommentsInfo {
		if val.IfAnonymous == true {
			val.CommentingUser = "****"
		}
		CommentsNew = append(CommentsNew, val)
	}
	//根据articleInfo中的IfAnonymous 字段是否为true,
	//然后复制到新切片，为true则不查询用户头像，将作者替换为：匿名用户
	//根据文章作者Author 字段在user表中查询头像数据
	var UserProfilePicture string
	var Author string
	if articleInfo.IfAnonymous == true {
		//设置匿名作者默认头像
		UserProfilePicture = "/static/imgs/user/default.png"
		Author = "匿名用户"

	} else {
		//查询用户头像
		var userInfo models.Users
		userName := articleInfo.Author
		db.Where("username = ?", userName).First(&userInfo)
		UserProfilePicture = userInfo.PicturePath
		//获取实际作者信息
		Author = articleInfo.Author
	}

	//返回数据
	ctx.HTML(200, "article/details.html", gin.H{
		"message":            "success",
		"articleID":          articleInfo.ID,
		"BlogTitle":          articleInfo.BlogTitle,
		"BlogContent":        articleInfo.BlogContent,
		"ReleaseDate":        articleInfo.ReleaseDate,
		"Author":             Author,
		"NumberOfViews":      articleInfo.NumberOfViews,
		"Likes":              articleInfo.Likes,
		"UserProfilePicture": UserProfilePicture,
		"currentUser":        currentUser.UserName,
		"comments":           CommentsNew,
	})
}

//WriteArticlePage 文章编辑界面
func WriteArticlePage(ctx *gin.Context) {
	db := common.GetDB()
	session := sessions.Default(ctx)
	currentUser := session.Get("currentUser").(UserInfo)
	//从数据库中获取tags
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
	//<option value="1">月季花</option>
	ctx.HTML(200, "article/write.html", gin.H{
		"message": "success",
		"Author":  currentUser.UserName,
		"options": optionHtml,
	})
}

//WriteArticle 写文章POST请求
func WriteArticle(ctx *gin.Context) {
	db := common.GetDB()
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
	var articleInfo article.WriteArticleForm
	//获取POST请求中的文章数据
	err = ctx.ShouldBind(&articleInfo)
	if err != nil {
		ctx.HTML(422, "article/write.html", gin.H{
			"msg":     "错误! " + err.Error(),
			"style":   "alert alert-dismissible alert-danger",
			"Author":  currentUser.UserName,
			"options": optionHtml,
		})
		return
	}
	//获取作者，标题,内容，是否匿名等参数
	Author := articleInfo.Author
	Title := articleInfo.BlogTitle
	Overview := articleInfo.BlogContentOverview
	Content := articleInfo.Content
	Anonymous := articleInfo.Anonymous
	Tag := articleInfo.Tag

	//判断数据长度是否符合要求
	if len(Author) > 300 {
		ctx.HTML(422, "article/write.html", gin.H{
			"msg":     "错误! 用户名长度范围最大为:100!",
			"style":   "alert alert-dismissible alert-danger",
			"Author":  currentUser.UserName,
			"options": optionHtml,
		})
		return
	}
	if len(Title) > 300 {
		ctx.HTML(422, "article/write.html", gin.H{
			"msg":     "错误! 标题长度范围最大为:100!",
			"style":   "alert alert-dismissible alert-danger",
			"Author":  currentUser.UserName,
			"options": optionHtml,
		})
		return
	}
	if len(Overview) > 300 {
		ctx.HTML(422, "article/write.html", gin.H{
			"msg":     "错误! 文章概述长度范围最大为:100!",
			"style":   "alert alert-dismissible alert-danger",
			"Author":  currentUser.UserName,
			"options": optionHtml,
		})
		return
	}
	if len(Content) > 16383 {
		ctx.HTML(422, "article/write.html", gin.H{
			"msg":     "错误! 内容长度范围最大为:16383!",
			"style":   "alert alert-dismissible alert-danger",
			"Author":  currentUser.UserName,
			"options": optionHtml,
		})
		return
	}

	//判断作者是否为currentUser一致
	if currentUser.UserName != Author {
		ctx.HTML(422, "article/write.html", gin.H{
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

	//判断是否为匿名发布
	var AnonymousBool bool
	if Anonymous == "on" {
		AnonymousBool = true
	} else {
		AnonymousBool = false
	}
	//添加到数据库
	AddArticleInfo := models.Articles{
		ReleaseDate:         time.Now().Format("2006-01-02 15:04:05"),
		Author:              Author,
		BlogTitle:           Title,
		BlogContentOverview: Overview,
		BlogContent:         Content,
		Likes:               0,
		Comments:            0,
		NumberOfViews:       0,
		IfAnonymous:         AnonymousBool,
		Tag:                 Tag}
	db.Create(&AddArticleInfo)

	//返回HTML信息
	ctx.HTML(200, "article/write.html", gin.H{
		"msg":     "success:发布成功!可以到主页查看或者搜索查看文章!",
		"style":   "alert alert-success alert-dismissable",
		"Author":  currentUser.UserName,
		"options": optionHtml,
	})
}

// ReceivePicture 接收富文本编辑器中的图片文件
func ReceivePicture(ctx *gin.Context) {
	//db := common.GetDB()
	//获取图片文件
	file, err := ctx.FormFile("file")
	if err != nil {
		fmt.Println("接收文件错误:" + err.Error())
		return
	}
	//保存到指定文件
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ".png"
	curDir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前工作目录失败: " + err.Error())
		return
	}
	filePath := "/static/imgs/article/" + fileName
	err = ctx.SaveUploadedFile(file, curDir+filePath)
	if err != nil {
		fmt.Println("保存图片文件错误：" + err.Error())
		return
	}
	//接收成功之后将图片链接返回给前端，summernote中显示
	ctx.JSON(200, gin.H{
		"url": filePath,
	})
}

//CommentingArticles 在文章详情页面给文章评论
func CommentingArticles(ctx *gin.Context) {
	db := common.GetDB()
	var commentInfo article.AddCommentsForm
	//获取评论参数
	err := ctx.ShouldBind(&commentInfo)
	if err != nil {
		ctx.Redirect(http.StatusMovedPermanently, "/article/details")
		fmt.Println("comments表单错误:" + err.Error())
		return
	}
	//获取用户名,文章ID,评论信息，是否匿名
	username := commentInfo.UserName
	articleID := commentInfo.ArticleID
	content := commentInfo.Content
	anonymous := commentInfo.Anonymous

	url := "/article/details?id=" + articleID
	//判断数据长度 5<username<=20 5<password<=20 phone = 11
	if len(content) > 600 {
		fmt.Println("评论信息长度最大为:200!")
		ctx.Redirect(http.StatusMovedPermanently, url)
		return
	}
	//将articleID由string转换为 int
	articleIDint, err := strconv.Atoi(articleID)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		return
	}
	//根据查询文章是否存在
	var articleData models.Articles
	db.First(&articleData, articleIDint)
	if articleData.ID == 0 {
		fmt.Println("文章不存在!")
		ctx.Redirect(http.StatusMovedPermanently, url)
		return
	}
	//根据匿名字段判断是否为匿名
	var anonymousBool bool
	if anonymous == "on" {
		anonymousBool = true
	} else {
		anonymousBool = false
	}

	//写入数据库
	AddCommentsInfo := models.Comments{
		ArticleID:      articleIDint,
		CommentingUser: username,
		Content:        content,
		IfAnonymous:    anonymousBool}
	db.Create(&AddCommentsInfo)
	//重定向
	ctx.Redirect(http.StatusMovedPermanently, url)

}

// ArticleAddLikes 给文章点赞
func ArticleAddLikes(ctx *gin.Context) {
	db := common.GetDB()
	articleID := ctx.Query("id")
	redirectUrl := ctx.Query("url")
	//fmt.Println("id: " + articleID + " " + "url: " + redirectUrl)
	if articleID == "" {
		fmt.Println("id 参数为空!")
		//ctx.HTML(404, "error/404.html", gin.H{
		//	"msg": "错误! id 参数为空!",
		//})
		ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		return
	}
	if redirectUrl == "" {
		fmt.Println("url 参数为空!")
		//ctx.HTML(404, "error/404.html", gin.H{
		//	"msg": "错误! url 参数为空!",
		//})
		ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		return
	}
	//将articleID 由string 转换为int
	articleIDint, err := strconv.Atoi(articleID)
	if err != nil {
		fmt.Println("参数错误:" + err.Error())
		//ctx.HTML(404, "error/404.html", gin.H{
		//	"msg": "错误!  参数错误!",
		//})
		ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		return
	}
	//根据ID查询文章是否存在
	var articleInfo models.Articles
	db.First(&articleInfo, articleIDint)
	if articleInfo.ID == 0 {
		fmt.Println("未查询到文章信息,无法查看!")
		//ctx.HTML(404, "error/404.html", gin.H{
		//	"msg": "错误! 文章不存在!",
		//})
		ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		return
	}
	//获取点赞数并更新
	LikesNum := articleInfo.Likes
	LikesNum++
	db.Model(&articleInfo).Updates(models.Articles{
		Likes: LikesNum,
	})
	//重定向到url(临时重定向)
	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)

}
