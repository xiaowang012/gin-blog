package main

import (
	"gin-blog/common"
	"gin-blog/form/user"
	"gin-blog/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Cors 解决跨域问题的func
//func Cors() gin.HandlerFunc {
//	return func(context *gin.Context) {
//		method := context.Request.Method
//		context.Header("Access-Control-Allow-Origin", "*")
//		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
//		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
//		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
//		context.Header("Access-Control-Allow-Credentials", "true")
//		if method == "OPTIONS" {
//			context.AbortWithStatus(http.StatusNoContent)
//		}
//		context.Next()
//	}
//}

func main() {
	db := common.InitDB()
	//定义默认的gin路由器
	router := gin.Default()
	//router.Use(Cors())
	//加载模板文件
	router.LoadHTMLGlob("template/*/*")
	router.StaticFS("/static", http.Dir("./static"))

	//注册页面
	router.GET("/register", func(ctx *gin.Context) {
		ctx.HTML(200, "register.html", gin.H{
			"message": "welcome!",
		})

	})

	//用户注册POST请求 /api/auth/register
	router.POST("/api/auth/register")

	//登录页面
	router.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(200, "login.html", gin.H{
			"message": "welcome!",
		})
	})

	//登录验证账号密码
	router.POST("/login", func(ctx *gin.Context) {
		var login user.LoginForm
		err := ctx.ShouldBind(&login)
		if err != nil {
			ctx.HTML(422, "login.html", gin.H{"message": "必须填写用户名和密码!"})
			return
		} else {
			username := login.Username
			password := login.Password

			//判断数据长度 5<username<=20 5<password<=20
			if len(username) <= 5 && len(username) > 20 {
				ctx.HTML(422, "login.html", gin.H{"message": "用户名长度范围为:5-20!"})
				return
			}
			if len(password) <= 5 && len(password) > 20 {
				ctx.HTML(422, "login.html", gin.H{"message": "密码长度范围为:5-20!"})
				return
			}
			//验证账号密码
			var user models.Users
			db.Where("username = ?", username).First(&user)

			if user.ID == 0 {
				ctx.HTML(422, "login.html", gin.H{"message": "用户不存在!"})
				return
			}
			if user.HashPassword == password {
				//重定向到index
				ctx.Redirect(http.StatusMovedPermanently, "/index")

			} else {
				ctx.HTML(403, "login.html", gin.H{
					"message": "登录失败!密码错误!",
				})
			}

		}

	})

	//修改密码
	router.POST("/api/auth/update_password", func(ctx *gin.Context) {
		var changepassword user.ChangePasswordForm
		err := ctx.ShouldBind(&changepassword)
		if err != nil {
			ctx.JSON(422, gin.H{
				"message": err.Error(),
			})
			return
		} else {
			username := changepassword.Username
			oldpassword := changepassword.OldPassword
			newpassword := changepassword.NewPassword

			//校验原密码
			//验证账号密码
			var user models.Users
			db.Where("username = ?", username).First(&user)
			//fmt.Println(user.ID, user.Username, user.Password, user.Phone)
			if user.ID == 0 {
				ctx.JSON(422, gin.H{"message": "用户不存在!"})
				return
			}
			if user.HashPassword == oldpassword {
				//原密码认证通过，修改为新密码
				db.Model(&user).Update("password", newpassword)
				ctx.JSON(200, gin.H{
					"message": "change password success",
					"userID":  user.ID,
				})
			} else {
				ctx.JSON(403, gin.H{
					"message": "change password failed ,old password error",
					"userID":  user.ID,
				})
			}
		}

	})

	//删除用户
	router.POST("/api/auth/delete", func(ctx *gin.Context) {
		var deleteuser user.DeleteUserForm
		err := ctx.ShouldBind(&deleteuser)
		if err != nil {
			ctx.JSON(422, gin.H{
				"message": err.Error(),
			})
			return
		} else {
			username := deleteuser.Username
			//查询该用户是否存在
			var user models.Users
			db.Where("username = ?", username).First(&user)
			if user.ID == 0 {
				ctx.JSON(422, gin.H{"message": "用户不存在!"})
				return
			} else {
				//将用户的 Active 字段改为false
				db.Model(&user).Update("Active", false)
				ctx.JSON(200, gin.H{
					"message": "delete user success",
					"userID":  user.ID,
					"active":  user.Active,
				})

			}

		}

	})

	//主页
	router.GET("/index", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"message": "欢迎来到xx网站!",
		})
	})
	router.Run("0.0.0.0:5001")
}
