package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"miniblog/config"
	"miniblog/model"
	"miniblog/util"
)

func GetLoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token:=c.GetHeader(config.DefaultConfig.JwtToken)
		if !util.IsEmptyString(token){
			claims,err:=util.GetClaims(token)
			if err==nil{
				simpleUser:=new(model.SimpleUser)
				err=json.Unmarshal([]byte(claims.Subject),simpleUser)
				if err==nil{
					c.Set(config.DefaultConfig.User,simpleUser)
				}
			}
		}
		c.Next()
	}
}
