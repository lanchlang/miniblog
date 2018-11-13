package middleware

import (
	"github.com/gin-gonic/gin"
)

func GetLoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		//session := sessions.Default(c)
		//user:=session.Get("user")
		//if v,ok:=user.(models.User);ok{
		//	c.Set("user",v)
		//}
		c.Next()
	}
}
