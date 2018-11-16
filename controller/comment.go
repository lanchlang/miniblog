package controller

import (
	"github.com/gin-gonic/gin"
	"miniblog/model"
	"net/http"
	"strconv"
	"time"
)

type CommentForm struct {
	BlogId       int64  `sql:"b_id" json:"blog_id"`
	Content      map[string]interface{} `form:"content" json:"content" binding:"required"`
	Type         int8  `form:"type" json:"type" binding:"required"`  //类型
	ReplyId       int64 `form:"reply_u_id" json:"reply_id"`
	ReplyUsername     string `form:"reply_u_name" json:"reply_username"`
}

//创建blog
func CreateComment(ctx *gin.Context){
	var form BlogForm
	// This will infer what binder to use depending on the content-type header.
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "数据错误"})
		return
	}
	hasLogin,err:=hasUserLogin(ctx)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if !hasLogin{
		ctx.JSON(http.StatusForbidden, gin.H{"error": "请先登录"})
		return
	}
	user,err:=getUser(ctx)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	blog:=model.Blog{
		UserId:user.Id,
		Username:user.Username,
		Avatar:user.Avatar,
		Title:form.Title,
		Intro:form.Intro,
		CategoryName:form.CategoryName,
		CategoryId:form.CategoryId,
		Tags:form.Tags,
		Cover:form.Cover,
		Content:form.Content,
		Type:form.Type,
		AccessLimit:form.AccessLimit,
		CreateDate:time.Now(),
	}
	err=model.Save(&blog)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	ctx.JSON(http.StatusOK, blog)
	return
}

//通过Id删除blog
func DeleteCommentById(ctx *gin.Context){
	idStr:=ctx.Param("id")
	id,err:=strconv.ParseInt(idStr,10,64)
	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"参数错误"})
		return
	}
	user,err:=getUser(ctx)
	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"参数错误"})
		return
	}
	err=new(model.Blog).Delete(id,user.Id)
	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"参数错误"})
		return
	}
	ctx.JSON(http.StatusOK,gin.H{"message":"您已经成功删除"})
	return
}