package controller

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/validator.v2"
	"miniblog/model"
	"net/http"
)

//用户名是否存在
//用户名长度在1到20之间
type UsernameExistForm struct {
	Username string `validate:"min=1,max=20,regexp=^[a-zA-Z0-9\u4E00-\u9FFF]*$"`
}
func UsernameExist(ctx *gin.Context)  {
	username:=ctx.PostForm("username")
     form:=UsernameExistForm{Username:username}
	if err := validator.Validate(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名长度只能在1到20之间,其值为字母、数字和中文"})
	}
     exist,err:=model.Exist(&model.User{},"username",username)
     if err!=nil{
		 ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		 return
	 }
     if exist{
		 ctx.JSON(http.StatusOK, gin.H{"exist":true,"message": "用户名已经存在"})
		 return
	 }
	ctx.JSON(http.StatusOK, gin.H{"exist":false,"message": "用户名不存在"})
	return
}

type EmailExistForm struct{
	Email string `validate:"min=5,max=30,regexp=^([A-Za-z0-9])+@([A-Za-z0-9])+\\.([A-Za-z]{2,4})$"`
}
//邮箱是否存在
func EmailExist(ctx *gin.Context)  {
	email:=ctx.PostForm("email")
	form:=EmailExistForm{Email:email}
	if err := validator.Validate(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的邮箱"})
	}
	exist,err:=model.Exist(&model.User{},"email",email)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if exist{
		ctx.JSON(http.StatusOK, gin.H{"exist":true,"message": "邮箱已经存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"exist":false,"message": "邮箱不存在"})
	return
}

//提供验证码：提供验证码图片及验证码，以及带过期时间不可更改的随机验证字符串（使用redis保存，以session-id+字符串为key，设置过期时间，随机生成验证码）
func ProvideCaptcha(ctx *gin.Context){

}
//使用邮箱注册
//邮箱注册：
//1：验证是否带有验证字符串，并且是否过期，如果验证不通过或已经过期，返回并告知原因
//2：再次检查用户名和邮箱是否已经存在，如果已经存在，返回并告知原因
//3：将用户信息插入待验证的表中
//4：向注册邮箱发送邮件（邮件中包含链接，链接中包含验证信息）
//5：将信息保存在邮箱验证表中（包括用户名，邮箱，验证信息，过期时间等）
//6：返回，并让其验证
func RegisterWithEmail(ctx *gin.Context){

}

//电话号码是否存在
func PhoneExist(ctx *gin.Context)  {

}

//通过手机号码发送验证码
func SendCaptchaToPhone(ctx *gin.Context)  {

}
//验证验证码和手机号码是否匹配或过期
func VerifyPhoneCaptcha(ctx *gin.Context){

}
//电话注册：
//1：检查用户名和电话号码是否存在，如果已经存在，返回并告知原因
//2：将用户信息存入用户表中（期间可以设置一些用户的基本信息）
//3：返回成功信息和用户识别标志
func RegisterWithPhone(ctx *gin.Context){

}

//请求邮箱重置密码：
//1：获取邮箱，并在查找是否存在，如果不存在，则返回并说明
//2：向邮箱中发送重置密码链接（链接中包含有验证信息），将邮箱，验证信息，过期时间保存到邮箱重置密码验证表中
//3：返回成功信息，并告知用户可以通过邮箱重置密码了
func RequestResetPasswordThroughEmail(ctx *gin.Context){

}

//验证邮箱重置密码链接：
//1：提取链接中信息，并在邮箱重置密码验证表查找并验证，如果验证失败：那么返回并告知失败原因
//2：验证成功：跳转到密码重置的界面（携带一次性验证信息）
func VerifyResetPasswordLinkThroughEmail(ctx *gin.Context){

}


//完成密码重置：
//1：验证一次性验证信息，如果失败，返回并告知原因
//2：重置密码，并返回信息，并告知用户
func ResetPassword(ctx *gin.Context){

}

//请求电话重置密码：
//1：获取电话，并在查找是否存在，如果不存在，则返回并说明
//3：返回成功信息（携带一次性验证信息和用户id），并告知用户可以重置密码了
func RequestResetPasswordThroughPhone(ctx *gin.Context){

}


//邮箱登录：
//0：检查用户是否已经登录，如果已经登录，则返回并告知原因
//1：验证是否带有验证字符串，并且是否过期，如果验证不通过或已经过期，返回并告知原因
//2：再次邮箱是否已经存在，如果不存在，返回并告知原因
//3：校验邮箱和密码如果不成功，返回结果并告知
//4：加载用户信息，并返回
func LoginWithEmailAndPassword(ctx *gin.Context){

}

//电话号码登录：
//0：检查用户是否已经登录，如果已经登录，则返回并告知原因
//1：验证是否带有验证字符串，并且是否过期，如果验证不通过或已经过期，返回并告知原因
//2：再次电话是否已经存在，如果不存在，返回并告知原因
//3：校验电话和密码如果不成功，返回结果并告知
//4：加载用户信息，并返回

func LoginWithPhoneAndPassword(ctx *gin.Context){

}

//手机验证码登录：
//1：获取手机号和验证码，进行验证，如果验证成功，如果手机不存在，那么返回，并告知其直接进入注册资料填写阶段
//2：通过手机号码，获取客户用户资料并返回
func LoginWithPhoneAndCaptcha(ctx *gin.Context){

}

//第三方登陆
//获取要登录哪个网站，与其授权服务器交互，
//然后跳转授权页面
func LoginWithOauth(ctx *gin.Context){

}


//用户授权后，第三方页面会跳转到指定地址（并附带信息）
//我方提取所需信息，如果用户尚未注册，那么跳转到指定页面，并让用户填写所需信息。如果已经注册，那么获取用户信息并返回。
func LoginCallback(ctx *gin.Context){

}