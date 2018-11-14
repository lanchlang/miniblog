package model

import "time"

type EmailRegisterUser struct{
	tableName struct{} `sql:"t_email_register_user"`
	Id int `sql:",pk" json:"id"`
	Username string `sql:"username" json:"username"`
	Email string `sql:"email" json:"email"`
	Password string `sql:"password" json:"password"`
	Expires time.Time `sql:"expires" json:"expires"`
	VerifyCode string  `sql:"verify_code" json:"verify_code"`
}
