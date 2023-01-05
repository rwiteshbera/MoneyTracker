package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/rwiteshbera/MoneyTracker/utils"
	"net/http"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.Request.Header.Get("token")

		if authToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "you must be logged in to the server (unauthorized).",
			})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(authToken)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			c.Abort()
		}

		c.Set("phoneNumber", claims.PhoneNumber)
		c.Next()
	}
}
