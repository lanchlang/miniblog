package config

import "miniblog/util"

type CommonConfig struct{
	DefaultListSize int `json:"default_list_size"`  //默认列表长度
	JwtToken        string `json:"jwt_token"`     //jwt token的名字
}
var DefaultConfig=new(CommonConfig)
func init() {
    err:=util.LoadConfig("config/common.json",DefaultConfig)
    if err!=nil{
    	panic(err)
	}
}
