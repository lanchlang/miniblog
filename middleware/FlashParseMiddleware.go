package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func FlashParse() gin.HandlerFunc {
	return func(c *gin.Context) {
		query:=c.Request.URL.Query()
	flash:=map[string]string{}
		for flashKey, values := range query {
			if strings.Index(flashKey,"flash_")>=0{
				key:= flashKey[6:]
				if len(values)>0{
					flash[key]=values[0]
				}
			}
		}
		c.Set("flash",flash)
		c.Next()
	}
}
