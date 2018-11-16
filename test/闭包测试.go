package main

import "miniblog/model"

func main() {
	err:=model.Get(buildInstance(),10)
	if err!=nil{
		println(err.Error())
	}
}

func buildInstance() interface{}{
	return new(model.User)
}