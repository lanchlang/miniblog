package model

import "time"

type EmailRegisterUser struct{
	tableName struct{} `sql:"t_email_register_user"`
	Id int `sql:",pk" json:"id" form:"id"`
	Username string `sql:"username" json:"username" form:"username"`
	Email string `sql:"email" json:"email" form:"email"`
	Password string `sql:"password" json:"password" form:"password"`
	Expires time.Time `sql:"expires" json:"expires" form:"expires"`
	VerifyCode string  `sql:"verify_code" json:"verify_code" form:"verify_code"`
}

func NewEmailRegisterUser()interface{}{
	return new(EmailRegisterUser)
}

func NewEmailRegisterUserList()interface{}{
	var list []EmailRegisterUser
	return &list
}