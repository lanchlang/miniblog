package main

import (
	"github.com/go-pg/pg"
	"miniblog/model"
	"miniblog/util"
)

func main() {
	var users []model.User
	err:=model.GetListByKey(&users,"username","ok")
	if err==pg.ErrNoRows{
		println(err.Error())
	}
	user:=model.User{
		Username:"hello",
		Password:util.Hash("kkkkkkkk"),
	}
	err=model.Save(&user)
	if err!=nil{
		println(err.Error())
	}
}
