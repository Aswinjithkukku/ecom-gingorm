package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/aswinjithkukku/ecom-gingorm/initializer"
	"github.com/aswinjithkukku/ecom-gingorm/models"
	"github.com/aswinjithkukku/ecom-gingorm/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
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

	go utils.SendEmailWithoutHTML([]string{user.Email}, "Successfull Registered!!", "Thankyou for being our partner. Let's make some collaborations")

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User created successfully",
	})
}

// User signin function.
func UserSignIn(c *gin.Context) {
	var body struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	// Validation body.
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errArray []string
		for _, e := range validationErrors {
			tag := e.Tag()
			field := e.Field()
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

	var user models.Users

	result := initializer.DB.Find(&user, "email=?", body.Email)
	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Sorry user with this email do not exist",
		})
		c.Abort()
		return
	}

	credentialCheck := user.CheckPassword(body.Password)
	if credentialCheck != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		c.Abort()
		return
	}

	//  Creatting JWT Token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       user.Id,
		"exp":       time.Now().Add(time.Hour * 24 * 30).Unix(),
		"name":      user.Name,
		"email":     user.Email,
		"isBlocked": user.BlockStatus,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create token",
		})
		c.Abort()
		return
	}

	// Setting token to cookie.
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user": gin.H{
			"name":      user.Name,
			"email":     user.Email,
			"isBlocked": user.BlockStatus,
			"token":     tokenString,
		},
	})

}

// Validatin User function.
func ValidateUser(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(models.Users)

	var userResponse = struct {
		Id          int
		Name        string
		Email       string
		BlockStatus bool
	}{
		Id:          userObj.Id,
		Name:        userObj.Name,
		Email:       userObj.Email,
		BlockStatus: userObj.BlockStatus,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    userResponse,
	})
}

// Adding User profile.
func AddUserProfile(c *gin.Context) {
	auth, _ := c.Get("user")

	var body struct {
		PhoneNumber int    `json:"phoneNumber"`
		Country     string `json:"country"`
		City        string `json:"city"`
		PinCode     int    `json:"pinCode"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	profile := models.Profile{
		PhoneNumber: body.PhoneNumber,
		Country:     body.Country,
		City:        body.City,
		PinCode:     body.PinCode,
		UserId:      auth.(models.Users).Id,
	}

	result := initializer.DB.Create(&profile)

	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot Creae Profile. Try again!",
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"profile": profile,
	})
}
