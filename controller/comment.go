package controller

import (
	"github.com/gin-gonic/gin"
	"miniblog/config"
	"miniblog/model"
	"miniblog/util"
	"net/http"
	"strconv"
	"time"
)

type CommentForm struct {
	BlogId       int64  `form:"b_id" json:"blog_id" binding:"required"`
	Content      map[string]interface{} `form:"content" json:"content" binding:"required"`
	Type         int8  `form:"type" json:"type" binding:"required"`  //类型
	ReplyId       int64 `form:"reply_u_id" json:"reply_id"`
	ReplyUsername     string `form:"reply_u_name" json:"reply_username"`
}

//创建comment
func CreateComment(ctx *gin.Context){
	var form CommentForm
	// This will infer what binder to use depending on the content-type header.
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "数据错误"})
		return
	}
	user,err:=getUser(ctx)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
    if user.Id<=0{
		ctx.JSON(http.StatusForbidden, gin.H{"error": "请先登录"})
		return
	}
	comment:=model.Comment{
		UserId:user.Id,
		Username:user.Username,
		BlogId:form.BlogId,
		Content:form.Content,
		Type:form.Type,
		ReplyId:form.ReplyId,
		ReplyUsername:form.ReplyUsername,
		CreateDate:time.Now(),
	}
	err=model.Save(&comment)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	ctx.JSON(http.StatusOK, comment)
	return
}

//通过Id删除comment
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
	if user.Id<=0{
		ctx.JSON(http.StatusForbidden,gin.H{"error":"请先登录"})
		return
	}
	//TODO:如果是管理员，也可以删除

	err=new(model.Blog).Delete(id,user.Id)
	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"参数错误"})
		return
	}
	ctx.JSON(http.StatusOK,gin.H{"message":"您已经成功删除"})
	return
}

func ListCommentByBlog(ctx *gin.Context){
	lastIdStr:=ctx.Query("last_id")
	lastId:=int64(0)
	//是否为空
	if !util.IsEmptyString(lastIdStr){
		var err error
		lastId,err=strconv.ParseInt(lastIdStr,10,64)
		if err!=nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
			return
		}
	}
	//根据blogId来获取
	idStr:=ctx.Param("id")
	//是否为空
	if !util.IsEmptyString(idStr){
		var err error
		id,err:=strconv.ParseInt(idStr,10,64)
		if err!=nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
			return
		}
		blogs,err:=new(model.Comment).GetBlogComment(id,lastId,config.DefaultConfig.DefaultListSize)
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
			return
		}
		ctx.JSON(http.StatusOK, blogs)
		return
	}else{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
			return
	}
}
