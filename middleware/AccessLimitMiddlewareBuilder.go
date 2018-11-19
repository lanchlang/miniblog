package middleware

import (
	"github.com/gin-gonic/gin"
	"miniblog/config"
	"miniblog/model"
	"net/http"
)

func BuildAccessLimitMiddleware(accessLevelRequired int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInterface,exist:=ctx.Get(config.DefaultConfig.User)
		if exist{
			if user,ok:=userInterface.(*model.SimpleUser);ok{
				if user.AccessLevel>=accessLevelRequired{
					ctx.Next()
					return
				}
			}
		}
		ctx.AbortWithStatusJSON(http.StatusForbidden,gin.H{"error":"您没有权限"})
		return
	}
}
