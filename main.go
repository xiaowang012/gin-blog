package main

import (
	"gin-blog/common"
	"gin-blog/route"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	//初始化DB
	common.InitDB()
	//定义默认的gin路由器
	router := gin.Default()
	router = route.Route(router)
	//加载模板文件
	router.LoadHTMLGlob("template/**/*")
	router.StaticFS("/static", http.Dir("./static"))
	//run
	router.Run("0.0.0.0:5001")

}
