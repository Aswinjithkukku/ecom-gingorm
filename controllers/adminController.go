package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/aswinjithkukku/ecom-gingorm/initializer"
	"github.com/aswinjithkukku/ecom-gingorm/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductStruct struct {
	ShortName          string  `json:"shortName"`
	LongName           string  `json:"longName"`
	Cost               int     `json:"cost"`
	Price              int     `json:"price"`
	FinalPrice         int     `json:"finalPrice"`
	IsDiscount         bool    `json:"isDiscount"`
	DiscountType       *string `json:"discountType"`
	DiscountPrice      *int    `json:"discountPrice"`
	Description        string  `json:"description"`
	Stock              int     `json:"stock"`
	DealerName         string  `json:"dealerName"`
	DealerPlace        string  `json:"dealerPlace"`
	HeroImage          string  `json:"heroImage"`
	ProductDestination string  `json:"productDestination"`
}

func AdminCreateProduct(c *gin.Context) {

	shortName := c.PostForm("shortName")
	longName := c.PostForm("longName")
	costString := c.PostForm("cost")
	cost, _ := strconv.Atoi(costString)

	priceString := c.PostForm("price")
	price, _ := strconv.Atoi(priceString)

	isDiscountString := c.PostForm("isDiscount")
	isDiscount := false
	if isDiscountString != "" {
		isDiscountParsed, _ := strconv.ParseBool(isDiscountString)
		isDiscount = isDiscountParsed
	}

	discountType := c.PostForm("discountType")
	discountPriceString := c.PostForm("discountPrice")
	discountPrice, _ := strconv.Atoi(discountPriceString)

	description := c.PostForm("description")
	stockString := c.PostForm("stock")
	stock, _ := strconv.Atoi(stockString)

	dealerName := c.PostForm("dealerName")
	dealerPlace := c.PostForm("dealerPlace")
	productDestination := c.PostForm("productDestination")

	var finalPrice int = price

	if isDiscount {
		if discountType == models.Percentage {
			finalPrice = price - ((discountPrice * price) / 100)
		} else if discountType == models.Flat {
			finalPrice = price - discountPrice
		}
	} else {
		discountType = ""
		discountPrice = 0
	}

	if finalPrice < cost {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The discount price cannot exceed cost",
		})
		c.Abort()
		return
	}
	// heroImage adding.
	heroImagePath, _ := c.FormFile("heroImage")
	extension := filepath.Ext(heroImagePath.Filename)
	heroImage := uuid.New().String() + extension
	c.SaveUploadedFile(heroImagePath, "./public/images"+heroImage)

	product := models.Products{
		ShortName:          shortName,
		LongName:           longName,
		Cost:               uint(cost),
		Price:              uint(price),
		FinalPrice:         uint(finalPrice),
		IsDiscount:         isDiscount,
		DiscountType:       &discountType,
		DiscountPrice:      &discountPrice,
		Description:        description,
		Stock:              uint(stock),
		DealerName:         dealerName,
		DealerPlace:        dealerPlace,
		ProductDestination: productDestination,
		HeroImage:          "/public/images" + heroImage,
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
	result := initializer.DB.Find(&product, "id = ?", productId)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "No data found",
		})
		c.Abort()
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": result.Error,
		})
		c.Abort()
		return
	}

	product.ShortName = body.ShortName
	product.LongName = body.LongName
	product.Cost = uint(body.Cost)
	product.Price = uint(body.Price)
	product.DiscountType = body.DiscountType
	product.DiscountPrice = body.DiscountPrice
	product.Description = body.Description
	product.Stock = uint(body.Stock)
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
	result := initializer.DB.First(&product, "id = ?", productId)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": result.Error,
		})
		c.Abort()
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "No data found",
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
	if result.RowsAffected == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "No data found",
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

	result := initializer.DB.Find(&product, "id = ?", productId)

	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": result.Error,
		})
		c.Abort()
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "No data found",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"product": product,
	})
}
