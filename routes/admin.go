package routes

import (
	"github.com/aswinjithkukku/ecom-gingorm/controllers"
	"github.com/gin-gonic/gin"
)

func AdminRoute(router *gin.Engine) {
	admin := router.Group("/api/admin")
	{
		admin.POST("/createproduct", controllers.AdminCreateProduct)
		admin.PATCH("/updateproduct/:productid", controllers.AdminUpdateProduct)
		admin.DELETE("/deleteproduct/:productid", controllers.AdminDeleteProduct)
		admin.GET("/products", controllers.AdminGetAllProducts)
		admin.GET("/product/:productid", controllers.AdminGetSingleProduct)
	}
}
