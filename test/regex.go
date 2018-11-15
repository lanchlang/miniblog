package main

import "regexp"

func main() {
	matched,err:=regexp.MatchString("^[a-zA-Z0-9\u4E00-\u9FFF]{1,20}$","王朗#")
	if err!=nil{
		println(err.Error())
	}
	println(matched)
}
