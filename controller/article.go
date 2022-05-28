package controller

import (
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
		ctx.HTML(200, "index/index_search.html", gin.H{
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
	ctx.HTML(200, "index/index_search.html", gin.H{
		"currentUser": userinfoNew.UserName,
		"msg":         fmt.Sprintf("查询成功! 查询到%d条数据", len(articles)),
		"style":       "alert alert-success alert-dismissable",
		"articles":    articles,
		"messages":    messagesNew,
		"currentPage": 1,
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
	//根据文章作者Author 字段在user表中查询头像数据
	var userInfo models.Users
	userName := articleInfo.Author
	db.Where("username = ?", userName).First(&userInfo)
	UserProfilePicture := userInfo.PicturePath
	//更新浏览量 +1
	numberOfViews := articleInfo.NumberOfViews
	numberOfViews++
	db.Model(&articleInfo).Updates(models.Articles{NumberOfViews: numberOfViews})
	//返回数据
	ctx.HTML(200, "article/details.html", gin.H{
		"message":            "success",
		"BlogTitle":          articleInfo.BlogTitle,
		"BlogContent":        articleInfo.BlogContent,
		"ReleaseDate":        articleInfo.ReleaseDate,
		"Author":             articleInfo.Author,
		"NumberOfViews":      articleInfo.NumberOfViews,
		"UserProfilePicture": UserProfilePicture,
		"currentUser":        currentUser.UserName,
		"comments":           CommentsInfo,
	})
}

//WriteArticlePage 文章编辑界面
func WriteArticlePage(ctx *gin.Context) {
	session := sessions.Default(ctx)
	currentUser := session.Get("currentUser").(UserInfo)
	////获取redis中的编辑数据
	//rdb := common.GetRedis()
	//var articleData ArticlesInfoStruct
	//val, err := rdb.Get(context.Background(), currentUser.UserName).Bytes()
	//err = json.Unmarshal(val, &articleData)
	//if err != nil {
	//	fmt.Println("解析articleData错误：" + err.Error())
	//}
	ctx.HTML(200, "article/write.html", gin.H{
		"message": "success",
		"Author":  currentUser.UserName,
		//"title":    articleData.Title,
		//"Overview": articleData.Overview,
		//"Content":  articleData.Content,
	})
}

//WriteArticle 写文章POST请求
func WriteArticle(ctx *gin.Context) {
	db := common.GetDB()
	session := sessions.Default(ctx)
	currentUser := session.Get("currentUser").(UserInfo)
	var articleInfo article.WriteArticleForm
	//获取POST请求中的文章数据
	err := ctx.ShouldBind(&articleInfo)
	if err != nil {
		ctx.HTML(422, "article/write.html", gin.H{
			"msg":   "错误! " + err.Error(),
			"style": "alert alert-dismissible alert-danger",
		})
		return
	}
	//获取作者，标题,内容，是否匿名等参数
	Author := articleInfo.Author
	Title := articleInfo.BlogTitle
	Overview := articleInfo.BlogContentOverview
	Content := articleInfo.Content
	Anonymous := articleInfo.Anonymous

	//判断数据长度是否符合要求
	if len(Author) > 300 {
		ctx.HTML(422, "article/write.html", gin.H{
			"msg":   "错误! 用户名长度范围最大为:100!",
			"style": "alert alert-dismissible alert-danger",
		})
		return
	}
	if len(Title) > 300 {
		ctx.HTML(422, "article/write.html", gin.H{
			"msg":   "错误! 标题长度范围最大为:100!",
			"style": "alert alert-dismissible alert-danger",
		})
		return
	}
	if len(Overview) > 300 {
		ctx.HTML(422, "article/write.html", gin.H{
			"msg":   "错误! 文章概述长度范围最大为:100!",
			"style": "alert alert-dismissible alert-danger",
		})
		return
	}
	if len(Content) > 30000 {
		ctx.HTML(422, "article/write.html", gin.H{
			"msg":   "错误! 内容长度范围最大为:10000!",
			"style": "alert alert-dismissible alert-danger",
		})
		return
	}
	////存进redis 过期时间设置为2小时
	//rdb := common.GetRedis()
	////定义文章编辑数据结构体的值
	//var articlesData ArticlesInfoStruct
	//articlesData.Title = Title
	//articlesData.Overview = Overview
	//articlesData.Content = Content
	////将articlesData 转换为byte
	//byteData, err := json.Marshal(articlesData)
	//if err != nil {
	//	fmt.Println("articlesData转换数据错误: " + err.Error())
	//}
	//err = rdb.Set(context.Background(), currentUser.UserName, byteData, time.Hour*2).Err()
	//if err != nil {
	//	println("SET articlesData错误!:" + err.Error())
	//}
	//判断作者是否为currentUser一致
	if currentUser.UserName != Author {
		ctx.HTML(422, "article/write.html", gin.H{
			"msg":      "failed:发布失败!不允许修改其他的用户名提交!",
			"style":    "alert alert-dismissible alert-danger",
			"Author":   currentUser.UserName,
			"title":    Title,
			"Overview": Overview,
			"Content":  Content,
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
	AddArticleInfo := models.Articles{ReleaseDate: time.Now().Format("2006-01-02 15:04:05"),
		Author: Author, BlogTitle: Title, BlogContentOverview: Overview, BlogContent: Content,
		Likes: 0, Comments: 0, NumberOfViews: 0, IfAnonymous: AnonymousBool}
	db.Create(&AddArticleInfo)
	//返回HTML信息
	ctx.HTML(200, "article/write.html", gin.H{
		"msg":    "success:发布成功!可以到主页查看或者搜索查看文章!",
		"style":  "alert alert-success alert-dismissable",
		"Author": currentUser.UserName,
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
