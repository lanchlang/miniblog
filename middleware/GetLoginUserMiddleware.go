package middleware

import (
	"github.com/gin-gonic/gin"
	"miniblog/config"
	"miniblog/util"
)

func GetLoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token:=c.GetHeader(config.DefaultConfig.JwtToken)
		if util.IsEmptyString(token){
			c.Next()
		}else{
            claims,err:=util.GetClaims(token)
            if err!=nil{
            	c.Next()
			}else{
				//claims.Subject
			}
		}
	}
}
