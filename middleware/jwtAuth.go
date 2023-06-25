package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/globus303/sportujspolu/utils"
)

type AuthKey string

// const authKey AuthKey = "authentication"

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := utils.ExtractToken(c)
		if tokenString != "yourmother" {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		c.Next()
	}
}
