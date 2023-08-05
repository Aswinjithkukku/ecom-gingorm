package controllers

import (
	"net/http"

	"github.com/aswinjithkukku/ecom-gingorm/initializer"
	"github.com/aswinjithkukku/ecom-gingorm/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// User SignUp function.
func UserSignUp(c *gin.Context) {
	var user models.Users

	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	// Validation of req.body.
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errArray []string
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()
			if tag == "required" {
				errArray = append(errArray, field+" "+tag)
			} else {
				errArray = append(errArray, field+" should be "+tag)
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errArray[0],
		})
		c.Abort()
		return
	}

	// Hash Passowrd.
	if err := user.HashPassword(user.Password); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	// Create User
	result := initializer.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to create user. Try again!!",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}
