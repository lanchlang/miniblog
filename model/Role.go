package model

type Role struct{
	tableName struct{} `sql:"t_role"`
	Id int `sql:"id,pk" json:"id" form:"id"`
	Name string `sql:"name,pk" json:"name" form:"name"`
	Auths []int `pg:"auths,array" json:"auths" form:"auths"`
}

func NewRole()interface{}{
	return new(Role)
}

func NewRoleList()interface{}{
	var list []Role
	return &list
}