package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"encoding/gob"
	"github.com/gin-contrib/sessions/memstore"
	"miniblog/middleware"
	"miniblog/controller"
	"miniblog/model"
)


func main() {
	//启动定时爬取任务
	//go cronCrawl()
	router := gin.Default()
	//store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	store:= memstore.NewStore([]byte("secret"))
	//limitHttpMethod:=middleware.LimitHttpMethod(http.MethodGet,http.MethodPost)
	router.Use(sessions.Sessions("mfsf", store))
	router.Use(middleware.GetLoginUser())
	api:=router.Group("/api/v1")
	{
		apiPosts := api.Group("/users")
		{
			//检查用户名，邮箱，电话是否存在
			apiPosts.GET("/username_exist",controller.UsernameExist)
			apiPosts.GET("/email_exist",controller.EmailExist)
			apiPosts.GET("/phone_exist",controller.PhoneExist)
		}
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
		//检测图形验证码
		api.POST("/verify_captcha",controller.VerifyCaptcha)
		//发送手机验证码
		api.POST("/send_captcha_to_phone",controller.SendCaptchaToPhone)
		//检测电话验证码
		api.POST("/verify_phone_captcha",controller.VerifyPhoneCaptcha)
	}
	//验证通过邮箱发送的重置密码的链接
	router.GET("/verify_reset_password_link_through_email",controller.VerifyResetPasswordLinkThroughEmail)
	//验证注册邮箱
	router.GET("/verify_register_through_email",controller.VerifyRegisterThroughEmail)
	//提供图形验证码
	router.GET("/provide_captcha",controller.ProvideCaptcha)

	router.Run()
}



//要在session中保存数据的时候需要注册,基础类型和基础类型的list不用注册
//否则返回错误：gob: type not registered for interface: XXXX
func init() {
	gob.Register(&model.User{})
	gob.Register([]*model.User{})
}