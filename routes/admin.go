package routes

import (
	"github.com/aswinjithkukku/ecom-gingorm/controllers"
	"github.com/gin-gonic/gin"
)

func AdminRoute(router *gin.Engine) {
	admin := router.Group("/api/admin")
	{
		admin.POST("/createproduct", controllers.AdminCreateProduct)
		admin.PATCH("/updateproduct/:slug", controllers.AdminUpdateProduct)
		admin.DELETE("/deleteproduct/:slug", controllers.AdminDeleteProduct)
		admin.GET("/products", controllers.AdminGetAllProducts)
		admin.GET("/product/:slug", controllers.AdminGetSingleProduct)
		admin.GET("/alluser", controllers.GetAllUsers)
		admin.GET("/user/:userid", controllers.GetOneUser)
	}
}
