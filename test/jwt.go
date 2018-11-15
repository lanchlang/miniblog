package main

import (
	"github.com/dgrijalva/jwt-go"
	"miniblog/util"
)

func main()  {
	tokenStr,err:=util.GenToken(60*10,"langziguilai@yeah.net")
	if err!=nil{
		println(err.Error())
		return
	}
	token,err:=util.GetToken(tokenStr)
	if err!=nil{
		println(err.Error())
		return
	}
	if claims,ok:= token.Claims.(*jwt.StandardClaims);ok{
	    println(claims)
	}else{
		return
	}
}