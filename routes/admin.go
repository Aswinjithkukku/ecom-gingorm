package routes

import (
	"github.com/aswinjithkukku/ecom-gingorm/controllers"
	"github.com/gin-gonic/gin"
)

func AdminRoute(router *gin.Engine) {
	admin := router.Group("/api/admin")
	{
		admin.POST("/createproduct", controllers.AdminCreateProduct)
		admin.GET("/ping", controllers.Ping)
	}
}
