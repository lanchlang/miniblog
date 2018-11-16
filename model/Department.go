package model

type Department struct{
	tableName struct{} `sql:"t_departments"`
	Id int `sql:",pk" json:"id" form:"id"`
	Name string `sql:"name" json:"name" form:"name"`
}

func NewDepartment()interface{}{
	return new(Department)
}

func NewDepartmentList()interface{}{
	var list []Department
	return &list
}