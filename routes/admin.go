package routes

import (
	"github.com/aswinjithkukku/ecom-gingorm/controllers"
	"github.com/gin-gonic/gin"
)

func AdminRoute(router *gin.Engine) {
	admin := router.Group("/admin")
	{
		admin.GET("/createproduct", controllers.CreateProduct)
	}
}
