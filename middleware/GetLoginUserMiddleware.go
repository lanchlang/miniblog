package middleware

import (
	"github.com/gin-gonic/gin"
	"miniblog/util"
)

func GetLoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token:=c.GetHeader("token")
		if util.IsEmptyString(token){
			c.Next()
		}else{

		}
	}
}
