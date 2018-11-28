package main

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"miniblog/admincontroller"
	"miniblog/config"
	"miniblog/controller"
	"miniblog/middleware"
	"miniblog/model"
)


func main() {
	//启动定时爬取任务
	//go cronCrawl()
	router := gin.Default()
	//store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	//store:= memstore.NewStore([]byte("secret"))
	//limitHttpMethod:=middleware.LimitHttpMethod(http.MethodGet,http.MethodPost)
	//router.Use(sessions.Sessions("mfsf", store))
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.GetLoginUser())
	api:=router.Group("/api/v1")
	{
		//检查用户名，邮箱，电话是否存在
		api.GET("/username_exist",controller.UsernameExist)
		api.GET("/email_exist",controller.EmailExist)
		api.GET("/phone_exist",controller.PhoneExist)
		//用户通过email，phone，username，captcha登录
		api.POST("/login/email",controller.LoginWithEmailAndPassword)
		api.POST("/login/phone",controller.LoginWithPhoneAndPassword)
		api.POST("/login/username",controller.LoginWithUsernameAndPassword)
		api.POST("/login/captcha",controller.LoginWithPhoneAndCaptcha)
		//用户注册
		api.POST("/register/email",controller.RegisterWithEmail)
		api.POST("/register/phone",controller.RegisterWithPhone)
		//请求重置密码
		api.POST("/forget/password/email",controller.RequestResetPasswordThroughEmail)
		api.POST("/forget/password/phone",controller.RequestResetPasswordThroughPhone)
		//重置密码
		api.POST("/reset/password/",controller.ResetPassword)
		//提供图形验证码
		api.GET("/provide_captcha",controller.ProvideCaptcha)
		//检测图形验证码
		api.POST("/verify_captcha",controller.VerifyCaptcha)
		//发送手机验证码
		api.POST("/send_captcha_to_phone",controller.SendCaptchaToPhone)
		//检测电话验证码
		api.POST("/verify_phone_captcha",controller.VerifyPhoneCaptcha)
		//用户接口
		apiUsers := api.Group("/users")
		{
			//获取用户的blogs
			apiUsers.GET("/:id/blogs",controller.GetListBlogByUser)
		}
		//blog接口
		apiBlogs := api.Group("/blogs")
		{
			//增删改查，列表
			apiBlogs.GET("/:id",controller.GetBlogById)
			apiBlogs.GET("/:id/comments",controller.ListCommentByBlog)//获取博客的评论
			apiBlogs.GET("/:id/like",controller.LikeBlogById)
			apiBlogs.GET("/:id/unlike",controller.UnLikeBlogById)
			apiBlogs.DELETE("/:id",controller.DeleteBlogById)
			apiBlogs.PUT("/:id",controller.UpdateBlogById)
			apiBlogs.POST("/",controller.CreateBlog)
			apiBlogs.POST("/search",controller.SearchBlog)
			apiBlogs.POST("/ids",controller.GetListBlogByIds)
			apiBlogs.GET("/",controller.ListBlog) //通过参数获取列表
		}
        //comment接口
        apiComments:=api.Group("/comments")
        {
        	apiComments.POST("/",controller.CreateComment)
        	apiComments.DELETE("/:id",controller.DeleteCommentById)
		}
	}
	//管理接口
	adminApi:=router.Group("/api/admin/v1")
	{
		//部门接口
		adminApiDepartment:=adminApi.Group("/departments",middleware.BuildAccessLimitMiddleware(config.DefaultConfig.UserAdminAccessLevel))
		{
			//列表
			adminApiDepartment.GET("/",admincontroller.BuildList(model.NewDepartmentList))
			//单个
			adminApiDepartment.GET("/:id",admincontroller.BuildGet(model.NewDepartment))
			//创建
			adminApiDepartment.POST("/",admincontroller.BuildCreate(model.NewDepartment))
			//更新
			adminApiDepartment.PUT("/:id",admincontroller.BuildUpdate(model.NewDepartment,[]string{"name"}))
			//删除
			adminApiDepartment.DELETE("/:id",admincontroller.BuildDelete(model.NewDepartment))
		}
		////权限接口
		adminApiAuthority:=adminApi.Group("/auths",middleware.BuildAccessLimitMiddleware(config.DefaultConfig.UserAdminAccessLevel))
		{
			//列表
			adminApiAuthority.GET("/",admincontroller.BuildList(model.NewAuthorityList))
			//单个
			adminApiAuthority.GET("/:id",admincontroller.BuildGet(model.NewAuthority))
			//创建
			adminApiAuthority.POST("/",admincontroller.BuildCreate(model.NewAuthority))
			//更新
			adminApiAuthority.PUT("/:id",admincontroller.BuildUpdate(model.NewAuthority,[]string{"name","department_id"}))
			//删除
			adminApiAuthority.DELETE("/:id",admincontroller.BuildDelete(model.NewAuthority))
		}
		////角色接口
		adminApiRole:=adminApi.Group("/roles",middleware.BuildAccessLimitMiddleware(config.DefaultConfig.UserAdminAccessLevel))
		{
			//列表
			adminApiRole.GET("/",admincontroller.BuildList(model.NewRoleList))
			//单个
			adminApiRole.GET("/:id",admincontroller.BuildGet(model.NewRole))
			//创建
			adminApiRole.POST("/",admincontroller.BuildCreate(model.NewRole))
			//更新
			adminApiRole.PUT("/:id",admincontroller.BuildUpdate(model.NewRole,[]string{"name","auths"}))
			//删除
			adminApiRole.DELETE("/:id",admincontroller.BuildDelete(model.NewRole))
		}
		//用户接口
		adminApiUsers:=adminApi.Group("/users",middleware.BuildAccessLimitMiddleware(config.DefaultConfig.UserAdminAccessLevel))
		{
			//列表
			adminApiUsers.GET("/",admincontroller.BuildList(model.NewUserList))
			//单个
			adminApiUsers.GET("/:id",admincontroller.BuildGet(model.NewUser))
			//创建
			adminApiUsers.POST("/",admincontroller.BuildCreate(model.NewUser))
			//更新
			adminApiUsers.PUT("/:id",admincontroller.BuildUpdate(model.NewUser,[]string{"access_level"}))
			//删除
			adminApiUsers.DELETE("/:id",admincontroller.BuildDelete(model.NewRole))
		}
		//blog接口
		adminApiBlogs:=adminApi.Group("/blogs",middleware.BuildAccessLimitMiddleware(config.DefaultConfig.BlogAdminAccessLevel))
		{
			//列表
			adminApiBlogs.GET("/",admincontroller.BuildList(model.NewBlogList))
			//单个
			adminApiBlogs.GET("/:id",admincontroller.BuildGet(model.NewBlog))
			//创建
			adminApiBlogs.POST("/",admincontroller.BuildCreate(model.NewBlog))
			//更新
			adminApiBlogs.PUT("/:id",admincontroller.BuildUpdate(model.NewBlog,[]string{"access_limit"}))
			//删除
			adminApiBlogs.DELETE("/:id",admincontroller.BuildDelete(model.NewBlog))
		}
		//comment接口
		adminApiComments:=adminApi.Group("/comments",middleware.BuildAccessLimitMiddleware(config.DefaultConfig.CommentAdminAccessLevel))
		{
			//列表
			adminApiComments.GET("/",admincontroller.BuildList(model.NewCommentList))
			//删除
			adminApiComments.DELETE("/:id",admincontroller.BuildDelete(model.NewBlog))
		}
	}
	//验证通过邮箱发送的重置密码的链接
	router.GET("/verify_reset_password_link_through_email",controller.VerifyResetPasswordLinkThroughEmail)
	//验证注册邮箱
	router.GET("/verify_register_through_email",controller.VerifyRegisterThroughEmail)


	router.Run()
}



//要在session中保存数据的时候需要注册,基础类型和基础类型的list不用注册
//否则返回错误：gob: type not registered for interface: XXXX
func init() {
	gob.Register(&model.User{})
	gob.Register([]*model.User{})
}