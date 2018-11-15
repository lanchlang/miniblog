package model

import (
	"time"
)

type PhoneCaptcha struct{
	tableName struct{} `sql:"t_phone_captcha"`
	Id int64 `sql:",pk" json:"id"`
	Phone string `sql:"phone" json:"phone"`
	Expires time.Time `sql:"expires" json:"expires"`
	VerifyCode string  `sql:"verify_code" json:"verify_code"`
}

//通过phone和verify_code获取验证表记录
func (phoneCaptcha *PhoneCaptcha) GetPhoneCaptcha(phone,captcha string)(error){
	session:=GetSession()
	defer session.Close()
	var captchas []PhoneCaptcha
	err:=session.Model(&captchas).Where("phone=? and verify_code=?",phone,captcha).Order("id desc").Select()
	if err!=nil{
		return err
	}
	//存在，则找最近的匹配
	if captchas!=nil && len(captchas)>0{
		phoneCaptcha.Phone=captchas[0].Phone
		phoneCaptcha.Id=captchas[0].Id
		phoneCaptcha.VerifyCode=captchas[0].VerifyCode
		phoneCaptcha.Expires=captchas[0].Expires
	}
	return nil
}