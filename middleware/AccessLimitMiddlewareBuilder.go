package middleware

import (
	"github.com/gin-gonic/gin"
)

func BuildAccessLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		//c.Set("flash",flash)
		c.Next()
	}
}
