package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-blog/common"
	"gin-blog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// ArticleList 文章列表页面
func ArticleList(ctx *gin.Context) {
	session := sessions.Default(ctx)
	currentUserInfo := session.Get("currentUser").(UserInfo)
	db := common.GetDB()
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

	//拼接HTML
	//按钮样式取随机值
	//按钮样式数组
	var tagStyle = [6]string{
		"btn btn-success",
		"btn btn-info",
		"btn btn-default",
		"btn btn-danger",
		"btn btn-warning",
		"btn btn-primary",
	}
	//传到前端的tags 切片
	var tag string
	var tagsHtml []string
	for _, v := range tags {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(6)
		tag = fmt.Sprintf("<a class=\"%s\" style=\"font-size: small;margin-left: 10px;margin-top: 10px;height: 30px;width: 60px;\" href=\"/article/list?tag=%s&page=1\">%s</a>",
			tagStyle[r], v.Tag, v.Tag)
		fmt.Println(tag)
		tagsHtml = append(tagsHtml, tag)
	}
	//获取tag参数
	tagParameter := ctx.Query("tag")
	pageNumber := ctx.Query("page")
	if tagParameter == "" {
		fmt.Println("tag 参数为空!")
		ctx.Redirect(http.StatusTemporaryRedirect, "/article/list?tag=Python&page=1")
		return
	}
	//将pageNumber转换为int
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Println("页码参数错误:" + err.Error())
		return
	}
	//根据tag查询article表
	var articleInfo []models.Articles
	db.Where("tag = ?", tagParameter).Limit(pageNumberInt * 10).Offset((pageNumberInt - 1) * 10).Order("created_at desc").Find(&articleInfo)
	//获取高赞文章，按点赞数量排序，获取前10条数据，并加入到redis中，过期时间设置2小时
	//渲染tags 标签
	var articleLikes []models.Articles
	valArticle, err := rdb.Get(context.Background(), "articleLikes").Bytes()
	if err != nil {
		fmt.Println("获取valArticle失败：" + err.Error())
		//从数据库中获取tags
		fmt.Println("从数据库中获取高赞文章-------------")
		db.Limit(10).Order("likes desc").Find(&articleLikes)
		byteData, err := json.Marshal(articleLikes)
		if err != nil {
			fmt.Println("tags 转换数据错误: " + err.Error())
		}
		err = rdb.Set(context.Background(), "articleLikes", byteData, time.Hour*2).Err()
		if err != nil {
			println("SET articleLikes错误!:" + err.Error())
		}
	} else {
		err = json.Unmarshal(valArticle, &articleLikes)
		if err != nil {
			fmt.Println("解析articleLikes失败-----")
			return
		}

	}
	//将高赞文章数据的title和id拼接成html a标签
	//<div><a href="/article/details?id=%d"><p style="color: black;margin-top: 10px;" >%s</p></a> </div>
	var articleLikesHtml []string
	var TagA string
	for _, value := range articleLikes {
		TagA = fmt.Sprintf("<div><a href=\"/article/details?id=%d\"><p style=\"color: black;margin-top: 10px;\" >%s</p></a> </div>",
			value.ID, value.BlogTitle)
		articleLikesHtml = append(articleLikesHtml, TagA)

	}
	ctx.HTML(200, "articleList/list.html", gin.H{
		"currentUser":      currentUserInfo.UserName,
		"tagsHtml":         tagsHtml,
		"articles":         articleInfo,
		"articleNum":       len(articleInfo),
		"currentPage":      pageNumberInt,
		"tag":              tagParameter,
		"articleLikesHtml": articleLikesHtml,
	})
}
