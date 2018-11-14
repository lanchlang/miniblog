package model

import (
	"github.com/go-pg/pg"
	"log"
	"time"
)

func onConnect(db *pg.DB) error {
	if db == nil {
		println("DataBase On Error")
	} else {
		println("DataBase On Connect")
	}
	return nil
}
func OnQueryProcessed(event *pg.QueryProcessedEvent) {
	query, err := event.FormattedQuery()
	if err != nil {
		panic(err)
	}

	log.Printf("%s %s", time.Since(event.StartTime), query)
}
func GetSession() *pg.DB {
	var session = pg.Connect(&pg.Options{
		Addr: "127.0.0.1",
		Database:  "miniblog",
		User:      "wanglang",
		Password:  "wanglang12345678",
		OnConnect: onConnect,
	})
	session.OnQueryProcessed(OnQueryProcessed)
	return session
}
func init() {
	//TEST connection
	session := pg.Connect(&pg.Options{
		Addr: "127.0.0.1",
		Database:  "miniblog",
		User:      "wanglang",
		Password:  "wanglang12345678",
		OnConnect: onConnect,
	})
	defer session.Close()
	//for _, model := range []interface{}{
	//	(*User)(nil),
	//	(*Oauth)(nil),
	//	(*Category)(nil),
	//	(*Blog)(nil),
	//	(*Verification)(nil),
	//	(*Comment)(nil),
	//} {
	//	err := session.CreateTable(model, &orm.CreateTableOptions{
	//		Temp:        false,
	//		IfNotExists: true,
	//	})
	//	if err != nil {
	//		print(err.Error())
	//		panic(err)
	//	}
	//}
}

//保存实例到数据库
//输入为一个实例model
//输出为error
func Save(model interface{}) error{
	session := GetSession()
	defer session.Close()
	return session.Insert(model)
}
//通过id从数据库中删除对应的记录
//输入为model(用于查找表)和id
//输出为error
func Delete(model interface{},id int64) error{
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Where("id=?", id).Delete()
	if err != nil {
		return err
	} else {
		return nil
	}
}

//通过id从数据库中获取记录
//输入为model(用于查找表)和id
//输出为error
//查找到的记录写入model中
func Get(model interface{},id int64) error {
	session := GetSession()
	defer session.Close()
	return session.Model(model).Where("id=?", id).Select()
}
//从数据库中获取列表
//输入为model(用于查找表)和limit,offset
//输出为error
//查找到的记录写入models中
func List(models []interface{},offset int,limit int) (error) {
	session := GetSession()
	defer session.Close()
	err := session.Model(&models).Offset(offset).Limit(limit).Select()
    return err
}
//更新某条记录
//输入为model，id，以及columns（需要更新的数组）
//输出为error
func Update(model interface{},id int,columns []string)(error){
	session := GetSession()
	defer session.Close()
	_,err:=session.Model(model).Column(columns...).Where("id=?",id).Update()
	return err
}
//获取某个表中的所有数据量
//输入为model
//输出为int(表中的数据量),error
func Count(model interface{}) (int,error) {
	session := GetSession()
	defer session.Close()
	count,err:=session.Model(model).Count()
	if err!=nil{
		return 0,err
	}
	return count,nil
}

//做简单模糊查询
//输入为models，查询关键字，查询的值，lastId,num
//输出为error
//查询结果集直接放到models中
func Search(key string,query string,lastId int,num int,models ...interface{}) error{
	session := GetSession()
	defer session.Close()
	err := session.Model(models).Where(key+` like ? and id<?`,"%"+query+"%",lastId).Limit(num).Select()
	return err
}

//检查表中的某些字段的值是否存在，例如用户名是否存在，email是否存在
//输入为model，表中的属性名称，需要匹配的值
func Exist(model interface{},property string,value interface{})(bool,error){
	session := GetSession()
	defer session.Close()
	count, err := session.Model(model).Where(property+"=?", value).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

//通过key值获取model
//输入为model，key,value
//输出为error
//查询结果集直接放到models中
func GetOneByKey(model interface{},key string,value interface{}) error{
	session := GetSession()
	defer session.Close()
	err := session.Model(model).Where(key+"=?",value).Select()
	return err
}
//通过key值获取models
//输入为models，key,value
//输出为error
//查询结果集直接放到models中
func GetListByKey(models interface{},key string,value interface{}) error{
	session := GetSession()
	defer session.Close()
	err := session.Model(models).Where(key+"=?",value).Select()
	return err
}