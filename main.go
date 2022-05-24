package main

import (
	"encoding/gob"
	"fmt"
	"gin-blog/common"
	"gin-blog/controller"
	"gin-blog/route"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	//session UserInfo 结构体
	gob.Register(controller.UserInfo{})
	//初始化DB
	common.InitDB()
	//初始化redis
	err := common.InitRedis()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//定义默认的gin路由器
	router := gin.Default()
	//定义session存储引擎redis
	store := cookie.NewStore([]byte("secretxaas121312xdff"))
	router.Use(sessions.Sessions("sessionID", store))
	router = route.Route(router)
	//加载模板文件
	router.LoadHTMLGlob("template/**/*")
	router.StaticFS("/static", http.Dir("./static"))
	//run
	router.Run("0.0.0.0:5001")

}
