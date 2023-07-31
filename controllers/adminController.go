package controllers

import (
	"fmt"
	"net/http"

	"github.com/aswinjithkukku/ecom-gingorm/initializer"
	"github.com/aswinjithkukku/ecom-gingorm/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductStruct struct {
	ShortName          string  `json:"shortName"`
	LongName           string  `json:"longName"`
	Cost               int     `json:"cost"`
	Price              int     `json:"price"`
	DiscountType       *string `json:"discountType"`
	DiscountPrice      *int    `json:"discountPrice"`
	Description        string  `json:"description"`
	Stock              int     `json:"stock"`
	DealerName         string  `json:"dealerName"`
	DealerPlace        string  `json:"dealerPlace"`
	ProductDestination string  `json:"productDestination"`
}

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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
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

func AdminUpdateProduct(c *gin.Context) {
	productId := c.Param("productid")

	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product!",
		})
		c.Abort()
		return
	}

	var body ProductStruct

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	var product models.Products

	result := initializer.DB.First(&product).Where("id = ?", productId)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": result.Error,
		})
		c.Abort()
		return
	}

	product.ShortName = body.ShortName
	product.LongName = body.LongName
	product.Cost = body.Cost
	product.Price = body.Price
	product.DiscountType = body.DiscountType
	product.DiscountPrice = body.DiscountPrice
	product.Description = body.Description
	product.Stock = body.Stock
	product.DealerName = body.DealerName
	product.DealerPlace = body.DealerPlace
	product.ProductDestination = body.ProductDestination

	initializer.DB.Save(&product)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"product": product,
	})

}

func AdminDeleteProduct(c *gin.Context) {
	productId := c.Param("productid")

	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Product",
		})
		c.Abort()
		return
	}

	var product models.Products
	result := initializer.DB.First(&product).Where("id = ?", productId)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": result.Error,
		})
		c.Abort()
		return
	}

	result = initializer.DB.Delete(&product)

	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": result.Error,
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, "Prodcut Deleted Successfully")

}

func AdminGetAllProducts(c *gin.Context) {
	var products []models.Products

	result := initializer.DB.Find(&products)

	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": result.Error,
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sucess":   true,
		"products": products,
	})
}

func AdminGetSingleProduct(c *gin.Context) {
	productId := c.Param("productid")

	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product",
		})
		c.Abort()
		return
	}

	var product models.Products

	result := initializer.DB.First(&product).Where("id = ?", productId)

	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": result.Error,
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"product": product,
	})
}

func Ping(c *gin.Context) {
	c.JSON(200, "pong")
}
