package controller

import (
	"github.com/gin-gonic/gin"
	"miniblog/model"
)

//检查用户是否登录
func hasUserLogin(ctx *gin.Context) (bool, error) {
	return false, nil
}

func getUser(ctx *gin.Context) (model.User,error){
	return model.User{},nil
}