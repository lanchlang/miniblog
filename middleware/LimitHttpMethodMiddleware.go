package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"errors"
)

func LimitHttpMethod(methods ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		isSupport:=false
		for _, method := range methods {
			if c.Request.Method == method {
				isSupport=true
				c.Next()
			}
		}
		if !isSupport{
			c.AbortWithError(http.StatusForbidden,errors.New("不该方法支持"))
		}

	}
}
