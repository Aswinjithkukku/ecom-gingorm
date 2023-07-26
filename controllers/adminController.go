package controllers

import (
	"fmt"
	"net/http"

	"github.com/aswinjithkukku/ecom-gingorm/initializer"
	"github.com/aswinjithkukku/ecom-gingorm/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func AdminCreateProduct(c *gin.Context) {

	var product models.Products

	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	validate := validator.New()
	if err := validate.Struct(product); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}

	result := initializer.DB.Create(&product)

	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add product",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"product": product,
	})

}

func Ping(c *gin.Context) {
	c.JSON(200, "pong")
}
