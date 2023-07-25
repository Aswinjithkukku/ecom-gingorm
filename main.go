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
}

func main() {

	r := gin.Default()

	routes.AdminRoute(r)

	r.Run(":" + os.Getenv("PORT"))
}