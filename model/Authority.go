package model

type Authority struct{
	tableName struct{} `sql:"t_authority"`
	Id int `sql:",pk" json:"id"`
	Name string `sql:"name" json:"name"`
	DepartmentId int `sql:"department_id" json:"department_id"`
}
