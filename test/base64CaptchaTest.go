package main

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
)

func main()  {
   demoCodeCaptchaCreate()
}

func demoCodeCaptchaCreate() {
	//config struct for digits
	//数字验证码配置
	var configD = base64Captcha.ConfigDigit{
		Height:     80,
		Width:      240,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 5,
	}
	//config struct for audio
	//声音验证码配置
	var configA = base64Captcha.ConfigAudio{
		CaptchaLen: 6,
		Language:   "zh",
	}
	//config struct for Character
	//字符,公式,验证码配置
	var configC = base64Captcha.ConfigCharacter{
		Height:             60,
		Width:              240,
		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode:               base64Captcha.CaptchaModeNumber,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsUseSimpleFont:    true,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         6,
	}
	//create a audio captcha.
	//GenerateCaptcha first parameter is empty string,so the package will generate a random uuid for you.
	idKeyA, capA := base64Captcha.GenerateCaptcha("", configA)
	//write to base64 string.
	//GenerateCaptcha first parameter is empty string,so the package will generate a random uuid for you.
	base64stringA := base64Captcha.CaptchaWriteToBase64Encoding(capA)
	//create a characters captcha.
	//GenerateCaptcha first parameter is empty string,so the package will generate a random uuid for you.
	idKeyC, capC := base64Captcha.GenerateCaptcha("", configC)
	//write to base64 string.
	base64stringC := base64Captcha.CaptchaWriteToBase64Encoding(capC)
	//create a digits captcha.
	idKeyD, capD := base64Captcha.GenerateCaptcha("", configD)
	//write to base64 string.
	base64stringD := base64Captcha.CaptchaWriteToBase64Encoding(capD)

	fmt.Println(idKeyA, base64stringA, "\n")
	fmt.Println(idKeyC, base64stringC, "\n")
	fmt.Println(idKeyD, base64stringD, "\n")
	base64Captcha.VerifyCaptcha(idKeyA,"")
}