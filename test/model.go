package main

import (
	"github.com/go-pg/pg"
	"miniblog/model"
)

func main() {
	var users []model.User
	err:=model.GetListByKey(&users,"username","ok")
	if err==pg.ErrNoRows{
		println(err.Error())
	}
}
