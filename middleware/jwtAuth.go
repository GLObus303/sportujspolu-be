package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/globus303/sportujspolu/constants"
	"github.com/globus303/sportujspolu/utils"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		userId, err := utils.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		c.Set(constants.UserID, userId)

		c.Next()
	}
}
