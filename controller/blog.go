package controller

import (
	"github.com/gin-gonic/gin"
	"miniblog/model"
	"net/http"
	"strconv"
)

func GetBlogById(ctx *gin.Context){
     idStr:=ctx.Param("id")
     id,err:=strconv.ParseInt(idStr,10,64)
     if err!=nil{
     	ctx.JSON(http.StatusBadRequest,gin.H{"error":"参数错误"})
     	return
	 }
     blog:=new(model.Blog)
     if err:=model.Get(blog,id);err!=nil{
		 ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		 return
	 }
     //获取部分评论
     comments,err:=new(model.Comment).GetBlogComment(id,0,20)
     if err!=nil{
		 ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		 return
	 }
	ctx.JSON(http.StatusOK, gin.H{"blog":blog,"comments":comments})
	return
}

//通过Id删除blog
func DeleteBlogById(ctx *gin.Context){
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
type BlogForm struct {
	Id           int64  `form:"id" json:"id"`
	Title        string `form:"title" json:"title"`
	Intro        string `form:"introduce" json:"introduce"`
	CategoryName string  `form:"category_name" json:"category_name"`
	CategoryId   int64   `form:"category_id" json:"category_id"`
	Tags         []string `form:"tags" json:"tags"`
	Cover        map[string]interface{} `form:"cover" json:"cover"`
	Content      map[string]interface{} `form:"content" json:"content"`
	Type         int8 `form:"type" json:"type"`   //类型
	AccessLimit  int  `form:"access_limit" json:"access_limit"`  //限制访问级别
}
func UpdateBlogById(ctx *gin.Context){
	idStr:=ctx.Param("id")
	id,err:=strconv.ParseInt(idStr,10,64)
	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"参数错误"})
		return
	}
	var form BlogForm
	// This will infer what binder to use depending on the content-type header.
	if err := ctx.ShouldBindJSON(&form); err != nil {
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
	if isCreator,err:=new(model.Blog).IsCreatorOfBlog(id,user.Id);err!=nil{
		if !isCreator{
			ctx.JSON(http.StatusForbidden, gin.H{"error": "您不是创建者，无法更新"})
			return
		}
	}
	blog:=model.Blog{
		Title:form.Title,
		Intro:form.Intro,
		CategoryName:form.CategoryName,
		CategoryId:form.CategoryId,
		Tags:form.Tags,
		Cover:form.Cover,
		Content:form.Content,
		Type:form.Type,
		AccessLimit:form.AccessLimit,
	}
    err=model.Update(&blog,id,[]string{"title","intro","category_name",
    "category_id","tags","cover","content","type","access_limit"})
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "您已经更新完成"})
	return
}

func CreateBlog(ctx *gin.Context){
	var form BlogForm
	// This will infer what binder to use depending on the content-type header.
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "数据错误"})
		return
	}
}

func SearchBlog(ctx *gin.Context){

}

func ListBlog(ctx *gin.Context){

}