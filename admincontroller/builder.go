package admincontroller

import (
	"github.com/gin-gonic/gin"
	"miniblog/config"
	"miniblog/model"
	"net/http"
	"strconv"
)

//获取
func BuildGet(instanceBuilder func() interface{}) func(ctx *gin.Context)  {
	return func (ctx *gin.Context){
		idStr:=ctx.Param("id")
		id,err:=strconv.ParseInt(idStr,10,64)
		if err!=nil{
			ctx.JSON(http.StatusBadRequest,gin.H{"error":"参数错误"})
			return
		}
        instance:=instanceBuilder()
        err=model.Get(instance,id)
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":"暂时不能服务"})
			return
		}
        ctx.JSON(http.StatusOK,instance)
		return
	}
}
//删除
func BuildDelete(instanceBuilder func() interface{}) func(ctx *gin.Context)  {
	return func (ctx *gin.Context){
		idStr:=ctx.Param("id")
		id,err:=strconv.ParseInt(idStr,10,64)
		if err!=nil{
			ctx.JSON(http.StatusBadRequest,gin.H{"error":"参数错误"})
			return
		}
		instance:=instanceBuilder()
		err=model.Delete(instance,id)
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":"暂时不能服务"})
			return
		}
		ctx.JSON(http.StatusOK,gin.H{"message":"删除成功"})
		return
	}
}
//获取列表
func BuildList(instanceBuilder func() interface{}) func(ctx *gin.Context)  {
	return func (ctx *gin.Context){
		offsetStr:=ctx.Query("offset")
		offset,err:=strconv.ParseInt(offsetStr,10,64)
		if err!=nil{
			ctx.JSON(http.StatusBadRequest,gin.H{"error":"参数错误"})
			return
		}
		instances:=instanceBuilder()
		err=model.List(instances,int(offset),config.Default_List_Size)
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":"暂时不能服务"})
			return
		}
		ctx.JSON(http.StatusOK,instances)
		return
	}
}

//更新
func BuildUpdate(instanceBuilder func() interface{}) func(ctx *gin.Context)  {
	return func (ctx *gin.Context){
		instance:=instanceBuilder()
		// This will infer what binder to use depending on the content-type header.
		if err := ctx.ShouldBindJSON(instance); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "数据错误"})
			return
		}
        err:=model.Update(instance)
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
			return
		}
		ctx.JSON(http.StatusOK,instance)
		return
	}
}

//创建
func BuildCreate(instanceBuilder func() interface{}) func(ctx *gin.Context)  {
	return func (ctx *gin.Context){
		instance:=instanceBuilder()
		// This will infer what binder to use depending on the content-type header.
		if err := ctx.ShouldBindJSON(instance); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "数据错误"})
			return
		}
		err:=model.Save(instance)
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
			return
		}
		ctx.JSON(http.StatusOK,instance)
		return
	}
}
