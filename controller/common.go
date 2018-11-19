package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"miniblog/config"
	"miniblog/model"
	"miniblog/util"
)

//检查用户是否登录
func hasUserLogin(ctx *gin.Context) (bool, error) {
	_,exist:=ctx.Get(config.DefaultConfig.User)
	return exist, nil
}

//获取用户
func getUser(ctx *gin.Context) (*model.SimpleUser,error){
	userInterface,exist:=ctx.Get(config.DefaultConfig.User)
	if exist{
		if user,ok:=userInterface.(*model.SimpleUser);ok{
			return user,nil
		}
	}
	return new(model.SimpleUser),nil
}

//设置jwtheader
func SetJwtHeader(ctx *gin.Context,user model.User) error{
	simpleUser:=model.SimpleUser{
		Id:user.Id,
		Username:user.Username,
		Email:user.Email,
		Phone:user.Phone,
		AccessLevel:user.AccessLevel,
	}
	bytes,err:=json.Marshal(&simpleUser)
	if err!=nil{
		return err
	}
	//长短时间为1周
	jwtStr,err:=util.GenToken(7*24*3600,string(bytes))
	ctx.Header(config.DefaultConfig.JwtToken,jwtStr)
	return nil
}