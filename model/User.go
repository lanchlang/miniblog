package model

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"miniblog/util"
	"time"
)

type SimpleUser struct{
	Id int64 `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	AccessLevel int `json:"access_level"`
}

type User struct {
	tableName struct{} `sql:"t_users"`
	Id          int64  `sql:"id,pk" json:"id" form:"id"`
	Username    string `sql:"username" json:"username" form:"username"`
	Password    string `sql:"password" json:"-"`
	Email       string `sql:"email" json:"email" form:"email"`
	Phone       string `sql:"phone" json:"phone" form:"phone"`
	CreateDate  time.Time `sql:"date_create" json:"create_date" form:"create_date"`
	LastLogin   time.Time `sql:"last_login" json:"last_login" form:"last_login"`
	Likes       []int64 `pg:"likes,array" json:"likes" form:"likes"`
	AccessLevel int     `sql:"access_level" json:"access_level" form:"access_level"`
}

func NewUser()interface{}{
	return new(User)
}

func NewUserList()interface{}{
	var list []User
	return &list
}

func (user *User) BeforeInsert(db orm.DB)  {
	if user.CreateDate.IsZero() {
		user.CreateDate = time.Now()
	}
	if user.LastLogin.IsZero() {
		user.LastLogin = time.Now()
	}
}

func (model *User) FindUserByUsernameAndPassword(username string, password string) error {
	session := GetSession()
	defer session.Close()
	err := session.Model(model).Where("username=? and password=?", username, util.Hash(password)).Select()
	if err ==pg.ErrNoRows {
		return nil
	}
	return err
}
func (model *User) FindUserByEmailAndPassword(email string, password string) error {
	session := GetSession()
	defer session.Close()
	err := session.Model(model).Where("email=? and password=?", email, util.Hash(password)).Select()

	return err
}
func (model *User) FindUserByPhoneAndPassword(phone string, password string) error {
	session := GetSession()
	defer session.Close()
	err := session.Model(model).Where("phone=? and password=?", phone, util.Hash(password)).Select()
	if err ==pg.ErrNoRows {
		return nil
	}
	return err
}
func (model *User) FindUserByEmail(email string) error {
	session := GetSession()
	defer session.Close()
	err := session.Model(model).Where("email=?", email).Select()
	if err ==pg.ErrNoRows {
		return nil
	}
	return err
}
func (model *User) SetEmailValid() error {
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Set("email_valid = true").Where("email = ?email").Update()
	if err != nil {
		print(err.Error())
	}
	return err
}
func (model *User) UpdatePassword(password string) error {
	session := GetSession()
	defer session.Close()
	model.Password = util.Hash(password)
	_, err := session.Model(model).Set("password = ?password").Where("id = ?id").Update()
	if err != nil {
		print(err.Error())
	}
	return err
}
func (model *User) Like() error {
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Set("like_cnt = like_cnt+1").Where("id = ?id").Update()
	if err != nil {
		print(err.Error())
	}
	return err
}
func (model *User) UnLike() error {
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Set("like_cnt = like_cnt-1").Where("id = ?id").Update()
	if err != nil {
		print(err.Error())
	}
	return err
}
func (model *User) CreateBlog() error {
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Set("blog_cnt = blog_cnt+1").Where("id = ?id").Update()
	if err != nil {
		print(err.Error())
	}
	return err
}
func (model *User) DeleteBlog() error {
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Set("blog_cnt = blog_cnt-1").Where("id = ?id").Update()
	if err != nil {
		print(err.Error())
	}
	return err
}
func (model *User) CreateCollector() error {
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Set("collection_cnt = collection_cnt+1").Where("id = ?id").Update()
	if err != nil {
		print(err.Error())
	}
	return err
}
func (model *User) DeleteCollector() error {
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Set("collection_cnt = collection_cnt-1").Where("id = ?id").Update()
	if err != nil {
		print(err.Error())
	}
	return err
}
func (model *User) FollowUser() error {
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Set("follow_user_cnt = follow_user_cnt+1").Where("id = ?id").Update()
	if err != nil {
		print(err.Error())
	}
	return err
}
func (model *User) UnFollowUser() error {
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Set("follow_user_cnt = follow_user_cnt-1").Where("id = ?id").Update()
	if err != nil {
		print(err.Error())
	}
	return err
}
func (model *User) SaveAvatar() error{
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Set("avatar = ?avatar").Where("id = ?id").Update()
	if err != nil {
		print(err.Error())
	}
	return err
}
