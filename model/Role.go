package model

type Role struct{
	tableName struct{} `sql:"t_role"`
	Id int `sql:"id,pk" json:"id"`
	Name string `sql:"name,pk" json:"name"`
	Auths []int `pg:"auths,array"`
}
