package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"encoding/gob"
	"miniblog/model"
	"github.com/gin-contrib/sessions/memstore"
	"miniblog/middleware"
	"miniblog/controller"
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
	//router.Static("/static", "./static")
	api:=router.Group("/api/v1")
	{
		apiPosts := api.Group("/users")
		{
			apiPosts.GET("/username_exist",controller.UsernameExist)
			apiPosts.GET("/email_exist",controller.EmailExist)
		}
		//user controller
		api.POST("/login",controller.ApiLogin)
		api.POST("/register",controller.ApiRegister)
		api.GET("/logout",controller.ApiLogout)
		api.POST("/reset/password",controller.ApiResetPassword)
		api.POST("/forget/password",controller.ApiForgetPassword)
	}

	router.Run()
}



//要在session中保存数据的时候需要注册,基础类型和基础类型的list不用注册
//否则返回错误：gob: type not registered for interface: XXXX
func init() {
	gob.Register(&models.User{})
	gob.Register([]*models.User{})
}