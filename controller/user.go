package controller

import (
	"github.com/gin-gonic/gin"
	"miniblog/config"
	"miniblog/model"
	"miniblog/util"
	"net/http"
	"strconv"
)

//获取某个用户的blogs
func GetListBlogByUser(ctx *gin.Context){
	userIdStr:=ctx.Param("id")
	userId,err:=strconv.ParseInt(userIdStr,10,64)
	if err!=nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}
	lastIdStr:=ctx.Query("last_id")
	lastId:=util.INT_64_MAX
	//是否为空
	if !util.IsEmptyString(lastIdStr){
		var err error
		lastId,err=strconv.ParseInt(lastIdStr,10,64)
		if err!=nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
			return
		}
	}
	user,err:=getUser(ctx)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if user.Id<=0{
		blogs,err:=new(model.Blog).GetPublicBlogByUser(userId,lastId,config.DefaultConfig.DefaultListSize)
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
			return
		}
		ctx.JSON(http.StatusOK, blogs)
		return
	}
	//TODO:如果是管理员，也可以获取所有的blogs
	if user.Id==userId{
		blogs,err:=new(model.Blog).GetAllBlogByUser(userId,lastId,config.DefaultConfig.DefaultListSize)
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
			return
		}
		ctx.JSON(http.StatusOK, blogs)
		return
	}

}
