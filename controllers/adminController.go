package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aswinjithkukku/ecom-gingorm/initializer"
	"github.com/aswinjithkukku/ecom-gingorm/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductStruct struct {
	ShortName          string `json:"shortName" validate:"required,min=3"`
	LongName           string `json:"longName"  validate:"required,min=3"`
	Cost               int    `json:"cost" validate:"required"`
	Price              int    `json:"price"  validate:"required"`
	FinalPrice         int    `json:"finalPrice"`
	IsDiscount         bool   `json:"isDiscount"`
	DiscountType       string `json:"discountType" `
	DiscountPrice      int    `json:"discountPrice"`
	Description        string `json:"description" validate:"required"`
	Stock              int    `json:"stock" validate:"required"`
	DealerName         string `json:"dealerName"`
	DealerPlace        string `json:"dealerPlace"`
	ProductDestination string `json:"productDestination"`
}

// Admin to create a product.
func AdminCreateProduct(c *gin.Context) {

	// Parsing value.
	cost, err := strconv.Atoi(c.PostForm("cost"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cost should be number type",
		})
	}
	price, err := strconv.Atoi(c.PostForm("price"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Price should be number type",
		})
	}
	isDiscountParsed, _ := strconv.ParseBool(c.PostForm("isDiscount"))

	discountPrice, err := strconv.Atoi(c.PostForm("discountPrice"))
	if err != nil && isDiscountParsed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Discount should be number type",
		})
	}

	stock, err := strconv.Atoi(c.PostForm("stock"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Stock should be number type",
		})
	}

	body := ProductStruct{
		ShortName:          c.PostForm("shortName"),
		LongName:           c.PostForm("longName"),
		Cost:               cost,
		Price:              price,
		IsDiscount:         isDiscountParsed,
		DiscountType:       c.PostForm("discountType"),
		DiscountPrice:      discountPrice,
		Description:        c.PostForm("description"),
		Stock:              stock,
		DealerName:         c.PostForm("dealerName"),
		DealerPlace:        c.PostForm("dealerPlace"),
		ProductDestination: c.PostForm("productDestination"),
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errArray []string
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()
			errArray = append(errArray, field+" "+tag)
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errArray[0],
		})
		c.Abort()
		return
	}

	var finalPrice int = price
	// Adding Discount Value
	if body.IsDiscount {
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
		IsDiscount:         body.IsDiscount,
		DiscountType:       nil,
		DiscountPrice:      nil,
		Description:        body.Description,
		Stock:              uint(stock),
		DealerName:         body.DealerName,
		DealerPlace:        body.DealerPlace,
		ProductDestination: body.ProductDestination,
		HeroImage:          "/public/" + folderName + "/" + heroImage,
	}

	if body.IsDiscount {
		product.DiscountType = &body.DiscountType
		product.DiscountPrice = &discountPrice
	}

	result := initializer.DB.Create(&product)

	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
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
	slug := c.Param("slug")

	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product!",
		})
		c.Abort()
		return
	}

	cost, err := strconv.Atoi(c.PostForm("cost"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cost should be number type",
		})
	}
	price, err := strconv.Atoi(c.PostForm("price"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Price should be number type",
		})
	}
	isDiscountParsed, _ := strconv.ParseBool(c.PostForm("isDiscount"))

	discountPrice, err := strconv.Atoi(c.PostForm("discountPrice"))
	if err != nil && isDiscountParsed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Discount should be number type",
		})
	}

	stock, err := strconv.Atoi(c.PostForm("stock"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Stock should be number type",
		})
	}

	body := ProductStruct{
		ShortName:          c.PostForm("shortName"),
		LongName:           c.PostForm("longName"),
		Cost:               cost,
		Price:              price,
		IsDiscount:         isDiscountParsed,
		DiscountType:       c.PostForm("discountType"),
		DiscountPrice:      discountPrice,
		Description:        c.PostForm("description"),
		Stock:              stock,
		DealerName:         c.PostForm("dealerName"),
		DealerPlace:        c.PostForm("dealerPlace"),
		ProductDestination: c.PostForm("productDestination"),
	}

	var product models.Products

	result := initializer.DB.Find(&product, "slug = ?", slug)

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
	dateTime := time.Now().Format("2006-01-02")
	dateValue, _ := time.Parse("2006-01-02", dateTime)
	folderName := dateValue.Format("2006-01-02")

	c.SaveUploadedFile(heroImagePath, "./public/"+folderName+"/"+heroImage)

	product.ShortName = body.ShortName
	product.LongName = body.LongName
	product.Cost = uint(cost)
	product.Price = uint(price)
	product.IsDiscount = body.IsDiscount
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
	slug := c.Param("slug")

	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Product",
		})
		c.Abort()
		return
	}

	var product models.Products
	result := initializer.DB.First(&product, "slug = ?", slug)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Sorry no data found!",
		})
		c.Abort()
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No data found!",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": result.Error.Error(),
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

	if result.RowsAffected == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "No data found",
		})
		c.Abort()
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": result.Error.Error(),
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
	slug := c.Param("slug")

	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product",
		})
		c.Abort()
		return
	}

	var product models.Products

	result := initializer.DB.Find(&product, "slug = ?", slug)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "No data found",
		})
		c.Abort()
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": result.Error.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"product": product,
	})
}

// Admin User controllers -------------------------------------------------------------------------------------------------.

// GetAllUsers function.
func GetAllUsers(c *gin.Context) {
	var users []models.Users

	result := initializer.DB.Preload("Profile").Find(&users)

	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot find the users",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"users":   users,
	})

}

// Get One User function.
func GetOneUser(c *gin.Context) {
	userId := c.Param("userid")

	var user models.Users

	result := initializer.DB.Preload("Profile").First(&user, "id=?", userId)

	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cannot find the provided user",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"user":    user,
	})
}
