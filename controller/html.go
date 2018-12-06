package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"miniblog/config"
	"miniblog/model"
	"miniblog/util"
	"net/http"
	"strconv"
)

func Main(ctx *gin.Context){
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
	//获取默认列表，根据创建日期（id）
	blogs,err:=new(model.Blog).GetPublicBlogByCreateDate(lastId,config.DefaultConfig.DefaultListSize)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	buffer := new(bytes.Buffer)
	//template.UserList(userList, buffer)
	ctx.Writer.Write(buffer.Bytes())
	ctx.JSON(http.StatusOK, blogs)
	return
}

//通过Id获取Blog
func Blog(ctx *gin.Context){
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
	//如果是非公开的
	if blog.AccessLimit>0{
		user,err:=getUser(ctx)
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
			return
		}

		//访问控制,有权限的人员可以访问
		if user.Id!=blog.UserId || user.AccessLevel<config.DefaultConfig.CommentAdminAccessLevel{
			ctx.JSON(http.StatusForbidden, gin.H{"error": "此文章暂时不对外开放"})
			return
		}
	}
	//visit + 1
	err=new(model.Blog).Visit(id)
	if err!=nil{
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

//受欢迎的blogs
func Popular(ctx *gin.Context){
	//根据popularOffset来获取
	offsetStr:=ctx.Query("offset")
	var offset int64=0
	//是否为空
	if !util.IsEmptyString(offsetStr){
		var err error
		offset,err=strconv.ParseInt(offsetStr,10,64)
		if err!=nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
			return
		}
	}
	blogs,err:=new(model.Blog).GetPublicBlogByPopular(int(offset),config.DefaultConfig.DefaultListSize)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	ctx.JSON(http.StatusOK, blogs)
}

//最受喜欢的blogs
func Favorite(ctx *gin.Context){
	var offset int64=0
	offsetStr:=ctx.Query("offset")
	//是否为空
	if !util.IsEmptyString(offsetStr){
		var err error
		offset,err=strconv.ParseInt(offsetStr,10,64)
		if err!=nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
			return
		}
	}
	blogs,err:=new(model.Blog).GetPublicBlogByLike(int(offset),config.DefaultConfig.DefaultListSize)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	ctx.JSON(http.StatusOK, blogs)
}

//根据tag来获取blogs
func Tag(ctx *gin.Context){
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
	//根据tag来获取
	tag:=ctx.Param("tag")
	//是否为空
	if !util.IsEmptyString(tag){
		var err error
		blogs,err:=new(model.Blog).GetPublicBlogByTag(tag,lastId,config.DefaultConfig.DefaultListSize)
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
			return
		}
		ctx.JSON(http.StatusOK, blogs)
	}
}
//根据category获取blogs
func Category(ctx *gin.Context){
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
	//根据categoryId来获取
	categoryIdStr:=ctx.Param("id")
	//是否为空
	if !util.IsEmptyString(categoryIdStr){
		var err error
		categoryId,err:=strconv.ParseInt(categoryIdStr,10,64)
		if err!=nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
			return
		}
		blogs,err:=new(model.Blog).GetPublicBlogByCategory(categoryId,lastId,config.DefaultConfig.DefaultListSize)
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
			return
		}
		ctx.JSON(http.StatusOK, blogs)
	}
}

//获取某个用户的blogs
func User(ctx *gin.Context){
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
	//如果是管理员或者是本人，可以获取所有的blogs
	if user.Id==userId || user.AccessLevel>=config.DefaultConfig.CommentAdminAccessLevel{
		blogs,err:=new(model.Blog).GetAllBlogByUser(userId,lastId,config.DefaultConfig.DefaultListSize)
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
			return
		}
		ctx.JSON(http.StatusOK, blogs)
		return
	}

}
