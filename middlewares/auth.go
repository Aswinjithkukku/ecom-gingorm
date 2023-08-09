package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aswinjithkukku/ecom-gingorm/initializer"
	"github.com/aswinjithkukku/ecom-gingorm/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func UserAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected Signing method : %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check for the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// find user with token sub
		var user models.Users
		initializer.DB.First(&user, claims["sub"])

		if user.Id == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// attach with req
		c.Set("user", user)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	c.Next()
}
