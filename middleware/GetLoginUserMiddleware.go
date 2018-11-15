package middleware

import (
	"github.com/gin-gonic/gin"
)

func GetLoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
