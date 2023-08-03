package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aswinjithkukku/ecom-gingorm/initializer"
	"github.com/aswinjithkukku/ecom-gingorm/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductStruct struct {
	ShortName          string `json:"shortName"`
	LongName           string `json:"longName"`
	Cost               string `json:"cost"`
	Price              string `json:"price"`
	FinalPrice         string `json:"finalPrice"`
	IsDiscount         string `json:"isDiscount"`
	DiscountType       string `json:"discountType"`
	DiscountPrice      string `json:"discountPrice"`
	Description        string `json:"description"`
	Stock              string `json:"stock"`
	DealerName         string `json:"dealerName"`
	DealerPlace        string `json:"dealerPlace"`
	ProductDestination string `json:"productDestination"`
}

// Admin to create a product.
func AdminCreateProduct(c *gin.Context) {

	body := ProductStruct{
		ShortName:          c.PostForm("shortName"),
		LongName:           c.PostForm("longName"),
		Cost:               c.PostForm("cost"),
		Price:              c.PostForm("price"),
		IsDiscount:         c.PostForm("isDiscount"),
		DiscountType:       c.PostForm("discountType"),
		DiscountPrice:      c.PostForm("discountPrice"),
		Description:        c.PostForm("description"),
		Stock:              c.PostForm("stock"),
		DealerName:         c.PostForm("dealerName"),
		DealerPlace:        c.PostForm("dealerPlace"),
		ProductDestination: c.PostForm("productDestination"),
	}

	// Parsing value.
	cost, _ := strconv.Atoi(body.Cost)
	price, _ := strconv.Atoi(body.Price)
	discountPrice, _ := strconv.Atoi(body.DiscountPrice)
	stock, _ := strconv.Atoi(body.Stock)

	isDiscount := false
	if body.IsDiscount != "" {
		isDiscountParsed, _ := strconv.ParseBool(body.IsDiscount)
		isDiscount = isDiscountParsed
	}

	var finalPrice int = price
	// Adding Discount Value
	if isDiscount {
		if body.DiscountType == models.Percentage {
			finalPrice = price - ((discountPrice * price) / 100)
		} else if body.DiscountType == models.Flat {
			finalPrice = price - discountPrice
		}
	}

	if finalPrice < cost {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The discount price cannot exceed cost",
		})
		c.Abort()
		return
	}
	// HeroImage adding.
	heroImagePath, err := c.FormFile("heroImage")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid image provided",
		})
		c.Abort()
		return
	}

	extension := filepath.Ext(heroImagePath.Filename)
	heroImage := uuid.New().String() + extension

	// Taking date for the folder.
	layout := "2006-01-02"
	dateTime := time.Now().Format("2006-01-02")
	dateValue, _ := time.Parse(layout, dateTime)
	folderName := dateValue.Format("2006-01-02")

	c.SaveUploadedFile(heroImagePath, "./public/"+folderName+"/"+heroImage)

	product := models.Products{
		ShortName:          body.ShortName,
		LongName:           body.LongName,
		Cost:               uint(cost),
		Price:              uint(price),
		FinalPrice:         uint(finalPrice),
		IsDiscount:         isDiscount,
		DiscountType:       nil,
		DiscountPrice:      nil,
		Description:        body.Description,
		Stock:              uint(stock),
		DealerName:         body.DealerName,
		DealerPlace:        body.DealerPlace,
		ProductDestination: body.ProductDestination,
		HeroImage:          "/public/" + folderName + "/" + heroImage,
	}

	if isDiscount {
		product.DiscountType = &body.DiscountType
		product.DiscountPrice = &discountPrice
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

// Admin to update product.
func AdminUpdateProduct(c *gin.Context) {
	productId := c.Param("productid")

	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product!",
		})
		c.Abort()
		return
	}

	body := ProductStruct{
		ShortName:          c.PostForm("shortName"),
		LongName:           c.PostForm("longName"),
		Cost:               c.PostForm("cost"),
		Price:              c.PostForm("price"),
		IsDiscount:         c.PostForm("isDiscount"),
		DiscountType:       c.PostForm("discountType"),
		DiscountPrice:      c.PostForm("discountPrice"),
		Description:        c.PostForm("description"),
		Stock:              c.PostForm("stock"),
		DealerName:         c.PostForm("dealerName"),
		DealerPlace:        c.PostForm("dealerPlace"),
		ProductDestination: c.PostForm("productDestination"),
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
	// HeroImage adding.
	heroImagePath, err := c.FormFile("heroImage")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid image provided",
		})
		c.Abort()
		return
	}

	extension := filepath.Ext(heroImagePath.Filename)
	heroImage := uuid.New().String() + extension

	// Taking date for the folder.
	layout := "2006-01-02"
	dateTime := time.Now().Format("2006-01-02")
	dateValue, _ := time.Parse(layout, dateTime)
	folderName := dateValue.Format("2006-01-02")

	c.SaveUploadedFile(heroImagePath, "./public/"+folderName+"/"+heroImage)

	isDiscount := false
	if body.IsDiscount != "" {
		isDiscountParsed, _ := strconv.ParseBool(body.IsDiscount)
		isDiscount = isDiscountParsed
	}

	cost, _ := strconv.Atoi(body.Cost)
	price, _ := strconv.Atoi(body.Price)
	discountPrice, _ := strconv.Atoi(body.DiscountPrice)
	stock, _ := strconv.Atoi(body.Stock)

	product.ShortName = body.ShortName
	product.LongName = body.LongName
	product.Cost = uint(cost)
	product.Price = uint(price)
	product.IsDiscount = isDiscount
	product.DiscountType = &body.DiscountType
	product.DiscountPrice = &discountPrice
	product.Description = body.Description
	product.Stock = uint(stock)
	product.DealerName = body.DealerName
	product.DealerPlace = body.DealerPlace
	product.ProductDestination = body.ProductDestination
	product.HeroImage = "/public/" + folderName + "/" + heroImage

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
