package controller

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/northbright/aliyun/message"
	"gopkg.in/validator.v2"
	"html/template"
	"miniblog/config"
	"miniblog/model"
	"miniblog/util"
	"net/http"
	"strings"
	"time"
)

//用户名是否存在
//用户名长度在1到20之间
type UsernameExistForm struct {
	Username string `validate:"min=1,max=20,regexp=^[a-zA-Z0-9\\u4E00-\\u9FFF]*$"`
}

func UsernameExist(ctx *gin.Context) {
	username := ctx.Query("username")
	form := UsernameExistForm{Username: username}
	if err := validator.Validate(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名长度只能在1到20之间,其值为字母、数字和中文"})
	}
	exist, err := model.Exist(&model.User{}, "username", username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if exist {
		ctx.JSON(http.StatusOK, gin.H{"exist": true, "message": "用户名已经存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"exist": false, "message": "用户名不存在"})
	return
}

type EmailExistForm struct {
	Email string `validate:"min=5,max=30,regexp=^([A-Za-z0-9])+@([A-Za-z0-9])+\\.([A-Za-z]{2\\,4})$"`
}

//邮箱是否存在
func EmailExist(ctx *gin.Context) {
	email := ctx.PostForm("email")
	form := EmailExistForm{Email: email}
	if err := validator.Validate(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的邮箱"})
	}
	exist, err := model.Exist(&model.User{}, "email", email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if exist {
		ctx.JSON(http.StatusOK, gin.H{"exist": true, "message": "邮箱已经存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"exist": false, "message": "邮箱不存在"})
	return
}

//提供验证码：https://github.com/mojocn/base64Captcha:支持多种样式,算术,数字,字母,混合模式,语音模式
//返回验证码的base64字符串和生成的验证码id
func ProvideCaptcha(ctx *gin.Context) {
	var config = base64Captcha.ConfigCharacter{
		Height: 60,
		Width:  240,
		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode:               base64Captcha.CaptchaModeNumberAlphabet,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsShowHollowLine:   true,
		IsShowNoiseDot:     true,
		IsShowNoiseText:    true,
		IsShowSlimeLine:    true,
		IsShowSineLine:     true,
		CaptchaLen:         6,
	}
	//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	idKey, cap := base64Captcha.GenerateCaptcha("", config)
	//以base64编码
	base64string := base64Captcha.CaptchaWriteToBase64Encoding(cap)
	ctx.JSON(http.StatusOK, gin.H{"id": idKey, "image": base64string})
}

//验证验证码：https://github.com/mojocn/base64Captcha:
//获取参数：id(验证码id)，value（验证码的值）
func VerifyCaptcha(ctx *gin.Context) {
	idkey := ctx.PostForm("id")
	verifyValue := ctx.PostForm("value")
	verifyResult := base64Captcha.VerifyCaptcha(idkey, verifyValue)
	if verifyResult {
		ctx.JSON(http.StatusOK, gin.H{"pass": true, "message": "验证通过"})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"pass": false, "message": "验证失败"})
	}
}

type RegisterWithEmailForm struct {
	Username string `validate:"min=1,max=20,regexp=^[a-zA-Z0-9\u4E00-\u9FFF]*$"`
	Email    string `validate:"min=5,max=30,regexp=^([A-Za-z0-9])+@([A-Za-z0-9])+\\.([A-Za-z]{2\\,4})$"`
	Password string `validate:"min=8,max=20,regexp=^[A-Za-z0-9]{8\\,20}"`
}

//使用邮箱注册
//邮箱注册：
//1：验证是否带有验证字符串，并且是否过期，如果验证不通过或已经过期，返回并告知原因
//2：再次检查用户名和邮箱是否已经存在，如果已经存在，返回并告知原因
//3：向注册邮箱发送邮件（邮件中包含链接，链接中包含验证信息）
//4：将信息保存在邮箱验证表中（包括用户名，邮箱，验证信息，过期时间等）
//5：返回，并让其验证
func RegisterWithEmail(ctx *gin.Context) {
	username := ctx.PostForm("username")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	form := RegisterWithEmailForm{Username: username,
		Email:    email,
		Password: password}
	if err := validator.Validate(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的参数"})
		return
	}
	//再次检查用户名和邮箱是否已经存在，如果已经存在，返回并告知原因
	exist, err := model.Exist(&model.User{}, "username", username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if exist {
		ctx.JSON(http.StatusOK, gin.H{"exist": true, "message": "用户名已经存在"})
		return
	}
	exist, err = model.Exist(&model.User{}, "email", email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if exist {
		ctx.JSON(http.StatusOK, gin.H{"exist": true, "message": "邮箱已经存在"})
		return
	}
	verifyCode := util.Hash(time.Now().String()) //生成验证码
	//向注册邮箱发送邮件（邮件中包含链接，链接中包含验证信息）
	mailHeader := "Email地址验证"
	mailBody,err := generateVerifyBodyContent("template/email_verify.html",verifyCode, email)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务，请稍后再试"})
		return
	}
	err = util.SendVerifyMail(email, mailHeader, mailBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务，请稍后再试"})
		return
	}
	//将用户信息插入待验证的表中
	user := &model.EmailRegisterUser{
		Username:   username,
		Email:      email,
		Password:   util.Hash(password),
		Expires:    time.Now().Add(168 * time.Hour),
		VerifyCode: verifyCode,
	}
	err = model.Save(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	//返回，并让其验证
	ctx.JSON(http.StatusOK, gin.H{"message": "非常感谢您的注册，请先验证您的邮箱，然后再登陆。提醒一下：您的邮箱验证的有效期为一周，过期将不再支持验证，谢谢！"})
	return
}

//生成邮箱验证内容
func generateVerifyBodyContent(templateName string,code string, email string) (string,error) {
	t,err:=template.ParseFiles(templateName)
	if err!=nil{
		return "",err
	}
	buf:=new(bytes.Buffer)
	err=t.Execute(buf, map[string]string{"code":code,"email":email})
	if err!=nil{
		return "",err
	}
	return buf.String(),nil
}

//验证邮箱
//1:通过验证码获取记录
//2:查看是否过期
//3:保存到主表中
//4:返回
func VerifyRegisterThroughEmail(ctx *gin.Context) {
	code := ctx.Query("code")
	var emailRegisterUser = new(model.EmailRegisterUser)
	//通过验证码获取记录
	err := model.GetOneByKey(emailRegisterUser, "verify_code", code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	//查看是否过期
	if time.Now().After(emailRegisterUser.Expires) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "邮箱链接已经过期，请重新注册，谢谢！"})
		return
	}
	//保存到主表中
	now := time.Now()
	user := model.User{
		Username:   emailRegisterUser.Username,
		Password:   emailRegisterUser.Password,
		Email:      emailRegisterUser.Email,
		CreateDate: now,
		LastLogin:  now,
		AccessLevel:config.DefaultConfig.DefaultUserAccessLevel,
	}
	err = model.Save(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	//设置jwt_header
    err=SetJwtHeader(ctx,user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
    //跳转到首页
    ctx.Redirect(http.StatusFound,"/")
}

type PhoneExistForm struct {
	Phone string `validate:"min=5,max=13,regexp=^1[34578]\\d{9}$"`
}

//电话号码是否存在
func PhoneExist(ctx *gin.Context) {
	phone := ctx.PostForm("phone")
	form := PhoneExistForm{Phone: phone}
	if err := validator.Validate(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的电话"})
		return
	}
	exist, err := model.Exist(&model.User{}, "phone", phone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if exist {
		ctx.JSON(http.StatusOK, gin.H{"exist": true, "message": "电话已经存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"exist": false, "message": "电话不存在"})
	return
}

//通过手机号码发送验证码
//TODO:通过服务商发送验证码:暂时使用阿里云：https://github.com/northbright/aliyun

type SendCaptchaToPhoneForm struct {
	Phone string `validate:"min=5,max=13,regexp=^1[34578]\\d{9}$"`
}

//发送手机验证码
func SendCaptchaToPhone(ctx *gin.Context) {
	//TODO:添加一些验证的数据，不能让机器发送
	phone := ctx.PostForm("phone")
	//校验
	form := SendCaptchaToPhoneForm{Phone: phone}
	if err := validator.Validate(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的电话"})
		return
	}
	//发送验证码
	verifyCode := util.GenerateRandomNumberString(4)
	err := sendMessageCaptcha(phone, verifyCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	//保存验证码
	phoneCaptcha := &model.PhoneCaptcha{
		Phone:      phone,
		VerifyCode: verifyCode,
		Expires:    time.Now().Add(15 * time.Minute), //15分钟之后失效
	}
	err = model.Save(phoneCaptcha)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "验证码已发送"})
}

//发送手机验证码
func sendMessageCaptcha(phone string, captcha string) error {
	var aliConfig config.AliConfig
	if err := config.LoadConfig("config/aliyun_message_config.json", &aliConfig); err != nil {
		return err
	}
	param := "{\"code\":\"" + captcha + "\"}"
	// Creates a new client.
	client := message.NewClient(aliConfig.AccessKeyID, aliConfig.AccessKeySecret)
	// Send SMS.
	ok, smsResp, err := client.SendSMS(
		[]string{phone},
		aliConfig.SMS.SignName,
		aliConfig.SMS.TemplateCode,
		param,
	)
	if err != nil {
		return err
	}
	if ok && smsResp.Code == "OK" {
		return nil
	}
	return errors.New("发生错误")
}

//验证验证码和手机号码是否匹配或过期
type VerifyPhoneCaptchaForm struct {
	Phone   string `validate:"min=5,max=13,regexp=^1[34578]\\d{9}$"`
	Captcha string `validate:"min=4,max=9,regexp=^\\d{4\\,9}$"`
}

//验证手机及其验证码
//TODO:记录错误到log
func VerifyPhoneCaptcha(ctx *gin.Context) {
	phone := ctx.PostForm("phone")
	captcha := ctx.PostForm("captcha")
	//校验
	form := VerifyPhoneCaptchaForm{Phone: phone, Captcha: captcha}
	if err := validator.Validate(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的电话和验证码"})
		return
	}
	//获取记录
	phoneCaptcha := new(model.PhoneCaptcha)
	err := phoneCaptcha.GetPhoneCaptcha(phone, captcha)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	//如果不存在这条记录
	if phoneCaptcha.Id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "验证失败"})
		return
	}
	//如果超时，返回说明
	if time.Now().After(phoneCaptcha.Expires) {
		//删除记录
		model.Delete(phoneCaptcha, phoneCaptcha.Id)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "验证码已超时"})
		return
	}
	//如果没有超时，则返回成功
	ctx.JSON(http.StatusOK, gin.H{"message": "验证成功"})
	return
}

//电话注册：
//1：检查用户名和电话号码是否存在，如果已经存在，返回并告知原因
//2：将用户信息存入用户表中（期间可以设置一些用户的基本信息）
//3：返回成功信息和用户识别标志
type RegisterWithPhoneForm struct {
	Username string `validate:"min=1,max=20,regexp=^[a-zA-Z0-9\u4E00-\u9FFF]*$"`
	Phone    string `validate:"min=5,max=13,regexp=^1[34578]\\d{9}$"`
	Password string `validate:"min=8,max=20,regexp=^[A-Za-z0-9]{8\\,20}"`
}

func RegisterWithPhone(ctx *gin.Context) {
	username := ctx.PostForm("username")
	phone := ctx.PostForm("phone")
	password := ctx.PostForm("password")
	form := RegisterWithPhoneForm{Username: username,
		Phone:    phone,
		Password: password}
	if err := validator.Validate(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的参数"})
		return
	}
	//再次检查用户名和电话是否已经存在，如果已经存在，返回并告知原因
	exist, err := model.Exist(&model.User{}, "username", username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if exist {
		ctx.JSON(http.StatusOK, gin.H{"exist": true, "message": "用户名已经存在"})
		return
	}
	exist, err = model.Exist(&model.User{}, "phone", phone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if exist {
		ctx.JSON(http.StatusOK, gin.H{"exist": true, "message": "电话已经存在，请通过电话重置密码"})
		return
	}
	//将用户信息插入表中
	user := model.User{
		Username:   username,
		Phone:      phone,
		Password:   util.Hash(password),
		CreateDate: time.Now(),
		LastLogin:  time.Now(),
		AccessLevel: config.DefaultConfig.BlogAdminAccessLevel,
	}
	err = model.Save(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	//
	//返回，并返回用户信息
	err=SetJwtHeader(ctx,user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "非常感谢您的注册！", "user": user})
	return
}

//请求邮箱重置密码：
//1：获取邮箱，并在查找是否存在，如果不存在，则返回并说明
//2：向邮箱中发送重置密码链接（链接中包含有验证信息），将邮箱，验证信息，过期时间保存到邮箱重置密码验证表中
//3：返回成功信息，并告知用户可以通过邮箱重置密码了
type RequestResetPasswordThroughEmailForm struct {
	Email string `validate:"min=5,max=30,regexp=^([A-Za-z0-9])+@([A-Za-z0-9])+\\.([A-Za-z]{2\\,4})$"`
}

func RequestResetPasswordThroughEmail(ctx *gin.Context) {
	email := ctx.PostForm("email")
	form := RequestResetPasswordThroughEmailForm{Email: email}
	if err := validator.Validate(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的邮箱"})
		return
	}
	//查找邮箱是否存在，如果不存在，则返回并说明
	exist, err := model.Exist(new(model.User), "email", email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if !exist {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "邮箱不存在，请先注册"})
		return
	}
	//向邮箱中发送重置密码链接（链接中包含有验证信息），将邮箱，验证信息，过期时间保存到邮箱重置密码验证表中
	verifyCode := util.Hash(time.Now().String()) //生成验证码
	mailHeader := "修改密码"
	mailBody,err := generateVerifyBodyContent("template/reset_password_through_email.html",verifyCode, email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务，请稍后再试"})
		return
	}
	err = util.SendVerifyMail(email, mailHeader, mailBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务，请稍后再试"})
		return
	}
	//将用户信息插入待验证的表中
	user := &model.EmailRegisterUser{
		Email:      email,
		Expires:    time.Now().Add(168 * time.Hour),
		VerifyCode: verifyCode,
	}
	err = model.Save(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	//返回，并让其验证
	ctx.JSON(http.StatusOK, gin.H{"message": "请登录您的邮箱进行密码重置，谢谢！"})
	return

}

//验证邮箱重置密码链接：
//1：提取链接中信息，并在邮箱重置密码验证表查找并验证，如果验证失败：那么返回并告知失败原因
//2：验证成功：跳转到密码重置的界面（携带一次性验证信息）
//TODO:跳转链接
func VerifyResetPasswordLinkThroughEmail(ctx *gin.Context) {
	code := ctx.Query("code")
	email := ctx.Query("email")
	var emailRegisterUser = new(model.EmailRegisterUser)
	//通过验证码获取记录
	err := model.GetOneByKey(emailRegisterUser, "verify_code", code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	//查看是否过期
	if time.Now().After(emailRegisterUser.Expires) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "邮箱链接已经过期，请重新使用邮箱重置，谢谢！"})
		return
	}
	//比较邮箱是否一致
	if strings.Compare(email, emailRegisterUser.Email) != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "邮箱有误"})
		return
	}
	tokenStr, err := util.GenToken(60*10, email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "您可以重置密码了", "token": tokenStr})
}

//完成密码重置：
//1：验证一次性验证信息，如果失败，返回并告知原因
//2：重置密码，并返回信息，并告知用户
type ResetPasswordForm struct {
	Password string `validate:"min=8,max=20,regexp=^[A-Za-z0-9]{8\\,20}"`
	TokenStr string
}

func ResetPassword(ctx *gin.Context) {
	password := ctx.PostForm("password")
	tokenStr := ctx.PostForm("token")
	form := ResetPasswordForm{Password: password, TokenStr: tokenStr}
	if err := validator.Validate(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的参数"})
		return
	}
	claims, err := util.GetClaims(tokenStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token参数有误"})
		return
	}
	subject := claims.Subject
	if isEmail, _ := util.EmailValid(subject); isEmail {
		user := model.User{}
		err = model.GetOneByKey(&user, "email", subject)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
			return
		}
		user.Password = util.Hash(password)
		err = model.UpdateColumns(&user, user.Id, []string{"password"})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
		return
	}
	if isPhone, _ := util.PhoneValid(subject); isPhone {
		user := model.User{}
		err = model.GetOneByKey(&user, "phone", subject)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
			return
		}
		user.Password = util.Hash(password)
		err = model.UpdateColumns(&user, user.Id, []string{"password"})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
		return
	}
	if isUsername, _ := util.PhoneValid(subject); isUsername {
		user := model.User{}
		err = model.GetOneByKey(&user, "username", subject)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
			return
		}
		user.Password = util.Hash(password)
		err = model.UpdateColumns(&user, user.Id, []string{"password"})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "不好意思，未找到用户"})
	return
}

//请求电话重置密码：
//1：获取电话，并在查找是否存在，如果不存在，则返回并说明
//3：返回成功信息（携带一次性验证信息和用户id），并告知用户可以重置密码了
type RequestResetPasswordThroughPhoneForm struct {
	Phone string `validate:"min=5,max=13,regexp=^1[34578]\\d{9}$"`
}

func RequestResetPasswordThroughPhone(ctx *gin.Context) {
	phone := ctx.PostForm("phone")
	form := RequestResetPasswordThroughPhoneForm{Phone: phone}
	if err := validator.Validate(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的电话"})
		return
	}
	//查找电话是否存在，如果不存在，则返回并说明
	exist, err := model.Exist(new(model.User), "phone", phone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if !exist {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "电话不存在，请先注册"})
		return
	}
	//发送验证码
	verifyCode := util.GenerateRandomNumberString(4)
	err = sendMessageCaptcha(phone, verifyCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	//保存验证码
	phoneCaptcha := &model.PhoneCaptcha{
		Phone:      phone,
		VerifyCode: verifyCode,
		Expires:    time.Now().Add(15 * time.Minute), //15分钟之后失效
	}
	err = model.Save(phoneCaptcha)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	//返回token
	tokenStr, err := util.GenToken(60*10, phone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "验证码已经发送", "token": tokenStr})
}

//邮箱登录：
//0：检查用户是否已经登录，如果已经登录，则返回并告知原因
//TODO：1：验证是否带有验证字符串，并且是否过期，如果验证不通过或已经过期，返回并告知原因
//3：校验邮箱和密码如果不成功，返回结果并告知
//4：加载用户信息，并返回
//TODO:加载额外信息
type LoginWithEmailAndPasswordForm struct {
	Email    string `validate:"min=5,max=30,regexp=^([A-Za-z0-9])+@([A-Za-z0-9])+\\.([A-Za-z]{2\\,4})$"`
	Password string `validate:"min=8,max=20,regexp=^[A-Za-z0-9]{8\\,20}"`
}

func LoginWithEmailAndPassword(ctx *gin.Context) {
	//检查用户是否已经登录，如果已经登录，则返回并告知原因
	hasLogin, err := hasUserLogin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if hasLogin {
		ctx.JSON(http.StatusOK, gin.H{"message": "您已经登录了，无需重新登录！"})
		return
	}
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	form := LoginWithEmailAndPasswordForm{
		Email:    email,
		Password: password,
	}
	if err := validator.Validate(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的邮箱和密码"})
		return
	}
	//查找邮箱是否存在，如果不存在，则返回并说明
	user := model.User{}
	err = user.FindUserByEmailAndPassword(email, password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if user.Id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "邮箱或者密码不正确"})
		return
	}
	//设置jwt_header
	err=SetJwtHeader(ctx,user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	//加载用户信息，并返回
	ctx.JSON(http.StatusOK, user)
	return
}



//电话号码登录：
//0：检查用户是否已经登录，如果已经登录，则返回并告知原因
//TODO：1：验证是否带有验证字符串，并且是否过期，如果验证不通过或已经过期，返回并告知原因
//2：再次电话是否已经存在，如果不存在，返回并告知原因
//3：校验电话和密码如果不成功，返回结果并告知
//4：加载用户信息，并返回
type LoginWithPhoneAndPasswordForm struct {
	Phone    string `validate:"min=5,max=13,regexp=^1[34578]\\d{9}$"`
	Password string `validate:"min=8,max=20,regexp=^[A-Za-z0-9]{8\\,20}"`
}

func LoginWithPhoneAndPassword(ctx *gin.Context) {
	//检查用户是否已经登录，如果已经登录，则返回并告知原因
	hasLogin, err := hasUserLogin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if hasLogin {
		ctx.JSON(http.StatusOK, gin.H{"message": "您已经登录了，无需重新登录！"})
		return
	}
	phone := ctx.PostForm("phone")
	password := ctx.PostForm("password")
	form := LoginWithPhoneAndPasswordForm{
		Phone:    phone,
		Password: password,
	}
	if err := validator.Validate(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的邮箱和密码"})
		return
	}
	//查找邮箱是否存在，如果不存在，则返回并说明
	user := model.User{}
	err = user.FindUserByPhoneAndPassword(phone, password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if user.Id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "电话或者密码不正确"})
		return
	}
	//设置jwt_header
	err=SetJwtHeader(ctx,user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	//加载用户信息，并返回
	ctx.JSON(http.StatusOK, user)
	return
}

//手机验证码登录：
//1：获取手机号和验证码，进行验证，如果验证成功，如果手机不存在，那么返回，并告知其直接进入注册资料填写阶段
//2：通过手机号码，获取客户用户资料并返回
type LoginWithPhoneAndCaptchaForm struct {
	Phone string `validate:"min=5,max=13,regexp=^1[34578]\\d{9}$"`
	Code  string `validate:"regexp=^\\d{4}$"`
}

func LoginWithPhoneAndCaptcha(ctx *gin.Context) {
	phone := ctx.PostForm("phone")
	code := ctx.PostForm("code")
	form := LoginWithPhoneAndCaptchaForm{
		Phone: phone,
		Code:  code,
	}
	if err := validator.Validate(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的电话和验证码"})
		return
	}
	captcha := &model.PhoneCaptcha{}
	if err := captcha.GetPhoneCaptcha(phone, code); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if captcha.Id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "验证失败"})
		return
	}
	if captcha.Expires.Before(time.Now()) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "验证码已超时"})
		return
	}
	user := model.User{}
	if err := model.GetOneByKey(&user, "phone", phone); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if user.Id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "电话不存在，请先注册"})
		return
	}
	//设置jwt_header
	err:=SetJwtHeader(ctx,user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	ctx.JSON(http.StatusOK, user)
	return
}

//电话号码登录：
//0：检查用户是否已经登录，如果已经登录，则返回并告知原因
//TODO：1：验证是否带有验证字符串，并且是否过期，如果验证不通过或已经过期，返回并告知原因
//3：校验用户名和密码如果不成功，返回结果并告知
//4：加载用户信息，并返回
type LoginWithUsernameAndPasswordForm struct {
	Username string `validate:"min=1,max=20,regexp=^[a-zA-Z0-9]{1\\,20}$"`
	Password string `validate:"min=8,max=20,regexp=^[A-Za-z0-9]{8\\,20}"`
}

func LoginWithUsernameAndPassword(ctx *gin.Context) {
	//检查用户是否已经登录，如果已经登录，则返回并告知原因
	hasLogin, err := hasUserLogin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if hasLogin {
		ctx.JSON(http.StatusOK, gin.H{"message": "您已经登录了，无需重新登录！"})
		return
	}
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	form := LoginWithUsernameAndPasswordForm{
		Username: username,
		Password: password,
	}
	if err := validator.Validate(form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的用户名"})
		return
	}
	//查找邮箱是否存在，如果不存在，则返回并说明
	user := model.User{}
	err = user.FindUserByUsernameAndPassword(username, password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	if user.Id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名或者密码不正确"})
		return
	}
	//设置jwt_header
	err=SetJwtHeader(ctx,user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "暂时不能服务"})
		return
	}
	//加载用户信息，并返回
	ctx.JSON(http.StatusOK, user)
	return
}

//第三方登陆
//获取要登录哪个网站，与其授权服务器交互，
//然后跳转授权页面
func LoginWithOauth(ctx *gin.Context) {

}

//用户授权后，第三方页面会跳转到指定地址（并附带信息）
//我方提取所需信息，如果用户尚未注册，那么跳转到指定页面，并让用户填写所需信息。如果已经注册，那么获取用户信息并返回。
func LoginCallback(ctx *gin.Context) {

}