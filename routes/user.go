package routes

import (
	"github.com/aswinjithkukku/ecom-gingorm/controllers"
	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine) {
	user := router.Group("/api/user")
	{
		user.POST("/signup", controllers.UserSignUp)
		user.POST("/signin", controllers.UserSignIn)
	}
}
