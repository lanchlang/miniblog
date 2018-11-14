package model

import "time"

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
	err:=session.Model(phoneCaptcha).Where("phone=? and verify_code=?",phone,captcha).Select()
	return err
}