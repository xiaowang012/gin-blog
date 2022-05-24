package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-blog/common"
	"gin-blog/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func IndexGET(ctx *gin.Context) {
	//获取redis连接
	rdb := common.GetRedis()
	db := common.GetDB()
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
	//定义接收mysql article，messages数据的slice
	var articles []models.Articles
	var messages []models.MessageBoard
	//判断redis中是否存在key articles，messages
	val, err := rdb.Get(context.Background(), "articles").Bytes()
	if err != nil {
		//fmt.Println("查询articles mysql数据")
		db.Limit(5).Offset(0).Find(&articles)
		byteData, err := json.Marshal(articles)
		if err != nil {
			fmt.Println("articles转换数据错误: " + err.Error())
		}
		err = rdb.Set(context.Background(), "articles", byteData, time.Minute*10).Err()
		if err != nil {
			println("SET articles错误!:" + err.Error())
		}
	} else {
		err = json.Unmarshal(val, &articles)
		if err != nil {
			fmt.Println("解析articles错误：" + err.Error())
		}
	}
	val1, err := rdb.Get(context.Background(), "messages").Bytes()
	if err != nil {
		db.Limit(5).Offset(0).Find(&messages)
		byteDataMessages, err := json.Marshal(messages)
		if err != nil {
			fmt.Println("messages转换数据错误: " + err.Error())
		}
		err = rdb.Set(context.Background(), "messages", byteDataMessages, time.Minute*10).Err()
		if err != nil {
			println("SET messages错误!:" + err.Error())
		}
	} else {
		err = json.Unmarshal(val1, &messages)
		if err != nil {
			fmt.Println("解析messages错误：" + err.Error())
		}
	}
	//返回数据到HTML
	ctx.HTML(200, "user/index.html", gin.H{
		"currentUser": userinfoNew.UserName,
		"msg":         userinfoNew.UserName + " ,欢迎来到Blog!",
		"style":       "alert alert-success alert-dismissable",
		"articles":    articles,
		"messages":    messages,
	})

}
