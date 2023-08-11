package main

import (
	"os"

	"github.com/aswinjithkukku/ecom-gingorm/initializer"
	"github.com/aswinjithkukku/ecom-gingorm/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.DatabaseConnection()
	initializer.SyncDatabase()
}

func main() {

	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	r.Static("/public", "./public")

	routes.AdminRoute(r)
	routes.UserRouter(r)

	r.Run(":" + os.Getenv("PORT"))
}
