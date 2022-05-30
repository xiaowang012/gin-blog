package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

//Logout 用户登出函数，删除session里面的user数据
func Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("currentUser")
	session.Save()
	ctx.Redirect(http.StatusTemporaryRedirect, "/login")
}
