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
	"html/template"
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
		fmt.Println("redis连接错误: " + err.Error())
		return
	}
	//定义默认的gin路由器
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		//模板中的两个变量相加，实现分页
		"add": func(x, y int) int {
			return x + y
		},
		//将返回到前端的字符串转换为html代码，如数据库字段为 <p>xx</p>,转换为前端的p标签
		"tran": func(code string) template.HTML {
			return template.HTML(code)
		},
	})
	//定义session存储引擎redis
	store := cookie.NewStore([]byte("13324@@3434341312admin"))
	router.Use(sessions.Sessions("sessionID", store))
	router = route.Route(router)
	//加载模板文件
	router.LoadHTMLGlob("template/**/*")
	router.StaticFS("/static", http.Dir("./static"))
	//run
	router.Run("0.0.0.0:5001")
}
