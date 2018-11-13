package model

type Department struct{
	tableName struct{} `sql:"t_departments"`
	Id int `sql:",pk" json:"id"`
	Name string `sql:"name" json:"name"`
}

