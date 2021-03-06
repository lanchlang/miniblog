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
//通过Id获取Blog
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
    if user.Id<=0{
		ctx.JSON(http.StatusForbidden,gin.H{"error":"请先登录"})
		return
	}
    //管理员删除博客
    if user.AccessLevel>=config.DefaultConfig.BlogAdminAccessLevel{
    	err=model.Delete(new(model.Blog),id)
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
			return
		}
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
	Title        string `form:"title" json:"title" binding:"required"`
	Intro        string `form:"introduce" json:"introduce" binding:"required"`
	CategoryName string  `form:"category_name" json:"category_name" binding:"required"`
	CategoryId   int64   `form:"category_id" json:"category_id" binding:"required"`
	Tags         []string `form:"tags" json:"tags" binding:"required"`
	Cover        map[string]interface{} `form:"cover" json:"cover" `
	Content      map[string]interface{} `form:"content" json:"content" binding:"required"`
	Type         int8 `form:"type" json:"type" binding:"required"`   //类型
	AccessLimit  int  `form:"access_limit" json:"access_limit" binding:"required"`  //限制访问级别
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
	user,err:=getUser(ctx)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if user.Id<=0{
		ctx.JSON(http.StatusForbidden, gin.H{"error": "请先登录"})
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
		Tags:form.Tags,
		Cover:form.Cover,
		Content:form.Content,
		Type:form.Type,
		AccessLimit:form.AccessLimit,
	}
    err=model.UpdateColumns(&blog,id,[]string{"title","intro","category_name",
    "category_id","tags","cover","content","type","access_limit"})
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "您已经更新完成"})
	return
}

//创建blog
func CreateBlog(ctx *gin.Context){
	var form BlogForm
	// This will infer what binder to use depending on the content-type header.
	if err := ctx.ShouldBindJSON(&form); err != nil {
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
	blog:=model.Blog{
		UserId:user.Id,
		Username:user.Username,
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

//搜索博客
//通过blog的title进行搜索
type SearchBlogForm struct {
	Query    string `validate:"min=1,max=20" form:"query" json:"query" binding:"required"`
	LastId   string `validate:"regexp=^[0-9]{0\\,20}$" form:"last_id" json:"last_id"`
}
func SearchBlog(ctx *gin.Context){
	var form SearchBlogForm
	// This will infer what binder to use depending on the content-type header.
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确搜索内容"})
		return
	}

	lastId:=util.INT_64_MAX
	//是否为空
    if !util.IsEmptyString(form.LastId){
    	var err error
		lastId,err=strconv.ParseInt(form.LastId,10,64)
		if err!=nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
			return
		}
	}
    //
    blogs,err:=new(model.Blog).SearchPublicBlog(form.Query,lastId,config.DefaultConfig.DefaultListSize)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	ctx.JSON(http.StatusOK, blogs)
	return
}

//按时间顺序逆序获取
func ListBlog(ctx *gin.Context){
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
	ctx.JSON(http.StatusOK, blogs)
	return
}


//受欢迎的blogs
func BlogsByPopular(ctx *gin.Context){
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
func BlogsByFavorite(ctx *gin.Context){
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
func BlogsByTag(ctx *gin.Context){
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
func BlogsByCategory(ctx *gin.Context){
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
//like post
func LikeBlogById(ctx *gin.Context){
	idStr:=ctx.Param("id")
	id,err:=strconv.ParseInt(idStr,10,64)
	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"参数错误"})
		return
	}
	user,err:=getUser(ctx)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
    if user.Id<=0{
		ctx.JSON(http.StatusForbidden, gin.H{"error": "请先登录或注册"})
		return
	}
	err=new(model.Blog).Like(user.Id,id)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "点赞成功"})
	return
}

//取消like post
func UnLikeBlogById(ctx *gin.Context){
	idStr:=ctx.Param("id")
	id,err:=strconv.ParseInt(idStr,10,64)
	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"参数错误"})
		return
	}
	user,err:=getUser(ctx)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if user.Id<=0{
		ctx.JSON(http.StatusForbidden, gin.H{"error": "请先登录或注册"})
		return
	}
	err=new(model.Blog).UnLike(user.Id,id)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "点赞成功"})
	return
}

//通过Id列表获取
type IdListForm struct{
	Ids []int64 `form:"ids" json:"ids" binding:"required"`
}
func GetListBlogByIds(ctx *gin.Context){
	var form IdListForm
	// This will infer what binder to use depending on the content-type header.
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ids数据错误"})
		return
	}
	//
	blogs,err:=new(model.Blog).QueryListByIds(form.Ids)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	user,err:=getUser(ctx)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	//未登录，只能返回公开的blog列表
	if user.Id<=0{
		var publicBlogs=[]model.Blog{}
		for _,blog:=range blogs{
			if blog.AccessLimit<=0{
				publicBlogs=append(publicBlogs,blog)
			}
		}
		ctx.JSON(http.StatusOK, publicBlogs)
		return
	}
	//T如果是管理员，则全部返回
    if user.AccessLevel>=config.DefaultConfig.CommentAdminAccessLevel{
    	ctx.JSON(http.StatusOK,blogs)
		return
	}
	//如果是登录用户，查看是否有权限，可以查看自己的blog
	var myBlogs=[]model.Blog{}
	for _,blog:=range blogs{
		if blog.AccessLimit<=user.AccessLevel{
			myBlogs=append(myBlogs,blog)
		}else{
			//如果是自己的blog
			if blog.UserId==user.Id{
				myBlogs=append(myBlogs,blog)
			}
		}
	}
	ctx.JSON(http.StatusOK, myBlogs)
	return
}
