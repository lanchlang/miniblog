package model

type Authority struct{
	tableName struct{} `sql:"t_authority"`
	Id int `sql:",pk" json:"id" form:"id"`
	Name string `sql:"name" json:"name" form:"name"`
	DepartmentId int `sql:"department_id" json:"department_id" form:"department_id"`
}

func NewAuthority()interface{}{
	return new(Authority)
}

func NewAuthorityList()interface{}{
	var list []Authority
	return &list
}