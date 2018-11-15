package main

import (
	"bytes"
	"html/template"
)

func main() {
	t,err:=template.ParseFiles("template/email_verify.html")
	if err!=nil{
		println(err.Error())
		return
	}
	buf:=new(bytes.Buffer)
	err=t.Execute(buf, map[string]string{"code":"http://www.bbbb.com","email":"lkk"})
	if err!=nil{
		println(err.Error())
		return
	}
	println(buf.String())
}
