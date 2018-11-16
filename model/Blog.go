package model

import (
	"github.com/go-pg/pg"
	"time"
	"github.com/go-pg/pg/orm"
)

type Blog struct {
	tableName struct{} `sql:"t_blog"`
	Id           int64 `sql:"id,pk" json:"id"`
	UserId       int64 `sql:"u_id" json:"user_id"`
	Username     string `sql:"u_name" json:"username"`
	Title        string `sql:"title" json:"title"`
	Avatar       map[string]interface{} `sql:"u_avatar" json:"avatar"`
	CreateDate   time.Time `sql:"date_create" json:"create_date"`
	CommentCnt   int `sql:"comment_cnt" json:"comment_cnt"`
	LikeCnt      int `sql:"like_cnt" json:"like_cnt"`
	Intro        string `sql:"intro" json:"introduce"`
	CategoryName string  `sql:"c_name" json:"category_name"`
	CategoryId   int64   `sql:"c_id" json:"category_id"`
	Tags         []string `pg:"tags,array"`
	Referer      string  `sql:"referer" json:"referer"`
	Cover        map[string]interface{} `sql:"cover" json:"cover"`
	Content      map[string]interface{} `sql:"content" json:"content"`
	UpdateDate   time.Time  `sql:"date_last_update" json:"update_date"`
	Type         int8 `sql:"type" json:"type"`   //类型
	AccessLimit  int  `sql:"access_limit" json:"access_limit"`  //限制访问级别
	ViewCnt     int   `sql:"view_cnt" json:"view_cnt"`//访问次数
	Deleted      int8  `sql:"deleted" json:"deleted"`   //deleted大于0，则表示已经删除
}

//获取在某个时间节点以后创建的blog
//输入为time.Time
//输出为int(数量),error
func (model *Blog) NewAddCountSince(time time.Time) (int,error) {
	session := GetSession()
	defer session.Close()
	count,err:=session.Model(model).Where("create_date>?",time).Count()
	if err!=nil{
		return 0,err
	}
	return count,nil
}
//检查blog的referer是否存在
//输入为referer
//输出为int(数量),error
func (model *Blog) RefererExist(referer string) (bool,error) {
	session := GetSession()
	defer session.Close()
	err:=session.Model(model).Column("referer").Where("referer=?", referer).Select()
	if err!=nil{
		return false,err
	}
	var exist bool
	if model.Id>0{
		exist=true
	}else{
		exist=false
	}
	return exist,nil
}
func (model *Blog) BeforeInsert(db orm.DB) error {
	if model.CreateDate.IsZero() {
		model.CreateDate = time.Now()
	}
	if model.UpdateDate.IsZero() {
		model.UpdateDate = time.Now()
	}
	model.LikeCnt = 0
	model.ViewCnt=0
	return nil
}

//用户删除自己的comment,deleted
//输入为comment的id和用户id
//输出为error
func (model *Blog) Delete(id int64, userId int64) error {
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Set("deleted=1").Where("id=? and u_id=?", id, userId).Update()
	return err
}
//用户更新自己的blog
//输入为blog的id,用户id以及需要更新的columns名称
//输出为error
//更新的值都放入到model中了
func (model *Blog) Update(id int64, userId int64,columns ...string) error {
	session := GetSession()
	defer session.Close()
	_,err:=session.Model(model).Column(columns...).Where("id=? and u_id=?",id,userId).Update()
	return err
}
//TODO 获取相似图片
func (model *Blog) GetSimlarBlog(page int64) ([]Blog, error) {
	numPerPage := 20
	offset := numPerPage * int(page-1)
	session := GetSession()
	defer session.Close()
	var blogs []Blog
	err := session.Model(&blogs).Where("id!=? and type=? and tags % ?", model.Id,model.Type,pg.Array(model.Tags)).Offset(offset).Limit(numPerPage).Select()
	//err := session.Model(&blogs).Where("id!=? and type=? and tags % ?", model.Id,model.Type,pg.Array(model.Tags)).OrderExpr("smlar(tags, ?, 'N.i / sqrt(N.a * N.b)' ) DESC",pg.Array(model.Tags)).Offset(offset).Limit(numPerPage).Select()
	if err != nil {
		print(err.Error())
		return nil, err
	} else {
		return blogs, nil
	}
}

//获取用户所有的blogs
//输入为user的id,上一个blog的id，获取的最大数量num
//输出为blog数组和error
//返回的blog数组是倒序排列的，最新的blog最先返回
func (model *Blog) GetAllBlogByUser(userId int64,lastBlogId int64, num int) ([]Blog, error) {

	session := GetSession()
	defer session.Close()
	var blogs []Blog
	err := session.Model(&blogs).Where("u_id=? and id<?", userId,lastBlogId).Limit(num).Select()
	if err != nil {
		return nil, err
	} else {
		return blogs, nil
	}
}

//获取用户所有用户公开的blogs
//输入为user的id,上一个blog的id，获取的最大数量num
//输出为blog数组和error
//返回的blog数组是倒序排列的，最新的blog最先返回
func (model *Blog) GetPublicBlogByUser(userId int64,lastBlogId int64, num int) ([]Blog, error) {
	session := GetSession()
	defer session.Close()
	var blogs []Blog
	err := session.Model(&blogs).Where("u_id=? and id<? and access_limit<=0", userId,lastBlogId).Limit(num).Select()
	if err != nil {
		print(err.Error())
		return nil, err
	} else {
		return blogs, nil
	}
}
//获取category里面所有的blogs
//输入为user的id,上一个blog的id，获取的最大数量num
//输出为blog数组和error
//返回的blog数组是倒序排列的，最新的blog最先返回
func (model *Blog) GetAllBlogByCategory(categoryId int64, lastBlogId int64, num int) ([]Blog, error) {
	session := GetSession()
	defer session.Close()
	var blogs []Blog
	err := session.Model(&blogs).Where("c_id=? and id<?", categoryId,lastBlogId).Order("id desc").Limit(num).Select()
	if err != nil {
		return nil, err
	} else {
		return blogs, nil
	}
}
//获取category里面所有公开的blogs
//输入为user的id,上一个blog的id，获取的最大数量num
//输出为blog数组和error
//返回的blog数组是倒序排列的，最新的blog最先返回
func (model *Blog) GetPublicBlogByCategory(categoryId int64, lastBlogId int64, num int) ([]Blog, error) {
	session := GetSession()
	defer session.Close()
	var blogs []Blog
	err := session.Model(&blogs).Where("c_id=? and id<? and access_limit<=0 ", categoryId,lastBlogId).Order("id desc").Limit(num).Select()
	if err != nil {
		return nil, err
	} else {
		return blogs, nil
	}
}

//获取所有带tag的blogs
//输入为user的id,上一个blog的id，获取的最大数量num
//输出为blog数组和error
//返回的blog数组是倒序排列的，最新的blog最先返回
func (model *Blog) GetAllBlogByTag(tag string, lastBlogId int64, num int) ([]Blog, error) {
	session := GetSession()
	defer session.Close()
	var blogs []Blog
	err := session.Model(&blogs).Where("?=ANY(tags) and id<?", tag,lastBlogId).Order("id desc").Limit(num).Select()
	if err != nil {
		print(err.Error())
		return nil, err
	} else {
		return blogs, nil
	}
}

//获取所有带tag的公开的blogs
//输入为user的id,上一个blog的id，获取的最大数量num
//输出为blog数组和error
//返回的blog数组是倒序排列的，最新的blog最先返回
func (model *Blog) GetPublicBlogByTag(tag string, lastBlogId int64, num int) ([]Blog, error) {
	session := GetSession()
	defer session.Close()
	var blogs []Blog
	err := session.Model(&blogs).Where("?=ANY(tags) and id<? and access_limit<=0", tag,lastBlogId).Order("id desc").Limit(num).Select()
	if err != nil {
		print(err.Error())
		return nil, err
	} else {
		return blogs, nil
	}
}

//获取热门的blogs
//输入为offset，获取的最大数量num
//输出为blog数组和error
//返回的blog数组是倒序排列的，最新的blog最先返回
func (model *Blog) GetAllBlogByPopular(offset int, num int) ([]Blog, error) {
	session := GetSession()
	defer session.Close()
	var blogs []Blog
	err := session.Model(&blogs).Order("visit_cnt DESC").Offset(offset).Limit(num).Select()
	if err != nil {
		print(err.Error())
		return nil, err
	} else {
		return blogs, nil
	}
}

//获取热门的公开的blogs
//输入为offset，获取的最大数量num
//输出为blog数组和error
//返回的blog数组是倒序排列的，最新的blog最先返回
func (model *Blog) GetPublicBlogByPopular(offset int, num int) ([]Blog, error) {
	session := GetSession()
	defer session.Close()
	var blogs []Blog
	err := session.Model(&blogs).Where("access_limit<=0").Order("visit_cnt DESC").Offset(offset).Limit(num).Select()
	if err != nil {
		print(err.Error())
		return nil, err
	} else {
		return blogs, nil
	}
}
//获取最受欢迎的blogs
//输入为offset，获取的最大数量num
//输出为blog数组和error
//返回的blog数组是倒序排列的，最新的blog最先返回
func (model *Blog) GetBlogByLike(offset int, num int) ([]Blog, error) {
	session := GetSession()
	defer session.Close()
	var blogs []Blog
	err := session.Model(&blogs).Order("like_cnt DESC").Offset(offset).Limit(num).Select()
	if err != nil {
		print(err.Error())
		return nil, err
	} else {
		return blogs, nil
	}
}
//获取最受欢迎的公开的blogs
//输入为offset，获取的最大数量num
//输出为blog数组和error
//返回的blog数组是倒序排列的，最新的blog最先返回
func (model *Blog) GetPublicBlogByLike(offset int, num int) ([]Blog, error) {
	session := GetSession()
	defer session.Close()
	var blogs []Blog
	err := session.Model(&blogs).Where("access_limit<=0").Order("like_cnt DESC").Offset(offset).Limit(num).Select()
	if err != nil {
		print(err.Error())
		return nil, err
	} else {
		return blogs, nil
	}
}

//根据创建时间(id)获取blogs
//输入为lastBlogId，获取的最大数量num
//输出为blog数组和error
//返回的blog数组是倒序排列的，最新的blog最先返回
func (blog *Blog) GetBlogByCreateDate(lastBlogId int64, num int) ([]Blog, error) {
	session := GetSession()
	defer session.Close()
	var blogs []Blog
	err := session.Model(&blogs).Where("id<?",lastBlogId).Order("id DESC").Limit(num).Select()
	if err != nil {
		print(err.Error())
		return nil, err
	} else {
		return blogs, nil
	}
}
//根据创建时间(id)获取公开的blogs
//输入为lastBlogId，获取的最大数量num
//输出为blog数组和error
//返回的blog数组是倒序排列的，最新的blog最先返回
func (blog *Blog) GetPublicBlogByCreateDate(lastBlogId int64, num int) ([]Blog, error) {
	session := GetSession()
	defer session.Close()
	var blogs []Blog
	err := session.Model(&blogs).Where("id<? and access_limit<=0",lastBlogId).Order("id DESC").Limit(num).Select()
	if err != nil {
		print(err.Error())
		return nil, err
	} else {
		return blogs, nil
	}
}

//通过传入id数组来获取blog
//输入为id的数组
//输出为blog数组和error
func (model *Blog) QueryListByIds(blogIds []int64) ([]Blog, error) {
	session := GetSession()
	defer session.Close()
	var blogs []Blog
	err := session.Model(&blogs).
		Where("id in (?)", pg.In(blogIds)).
		Select()
	if err != nil {
		print(err.Error())
		return nil, err
	} else {
		return blogs, nil
	}
}

//删除用户所有的blog
//输入为用户id
//输出为error
//仅仅更新一个字段，并没有从数据库中删除
func (blog *Blog) DeleteByUser(userId int64) (error) {
	session := GetSession()
	defer session.Close()
	_, err := session.Model(blog).Set("deleted=1").Where("u_id", userId).Update()
	return err
}
//搜索所有的Blog
//输入为query，lastBlogId,num
//输出为blog数组，error
//根据title来搜索
func (blog *Blog) SearchBlog(query string,lastBlogId int64,num int) ([]Blog, error) {
	session := GetSession()
	defer session.Close()
	var blogs []Blog
	err := session.Model(&blogs).Where(`title like ? and id<?`,"%"+query+"%",lastBlogId).
		Order("id desc").Limit(num).Select()
	if err != nil {
		print(err.Error())
		return nil, err
	} else {
		return blogs, nil
	}
}
//搜索所有的公开的Blog
//输入为query，lastBlogId,num
//输出为blog数组，error
//根据title来搜索
func (blog *Blog) SearchPublicBlog(query string,lastBlogId int64,num int) ([]Blog, error) {
	session := GetSession()
	defer session.Close()
	var blogs []Blog
	err := session.Model(&blogs).Where(`id<? and access_limit<=0 and title like ?`,"%"+query+"%",lastBlogId).
		Order("id desc").Limit(num).Select()
	if err != nil {
		print(err.Error())
		return nil, err
	} else {
		return blogs, nil
	}
}
//每次访问，都将增加一次view
func (model *Blog) Visit(id int64) ( error) {
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Set("view_cnt=view_cnt+1").
		Where("id = ?", id).Update()
	if err != nil {
		print(err.Error())
		return err
	}
	return nil

}
//点赞
func (model *Blog) Like(userId,id int64) ( error) {
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Set("like_cnt=like_cnt+1").
		Where("id = ?", id).Update()
	if err != nil {
		return err
	}
	_, err = session.Model(new(User)).Exec("UPDATE t_users SET likes = array_append(likes,?::BIGINT) WHERE id=? AND (?!=ALL(likes))", id, userId, id)
	if err != nil {
		return err
	} else {
		return nil
	}
	return nil
}
//取消点赞
func (model *Blog) UnLike(userId,id int64) ( error) {
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Set("like_cnt=like_cnt-1").
		Where("id = ?", id).Update()
	if err != nil {
		return err
	}
	_, err = session.Model(model).Exec("UPDATE t_users SET likes = array_remove(likes,?::BIGINT) WHERE id=?", id, userId)
	if err != nil {
		return err
	} else {
		return nil
	}
	return nil
}
//判断用户是否为blog创建者
func (model *Blog) IsCreatorOfBlog(id int64,userId int64)(bool,error){
	session := GetSession()
	defer session.Close()
	count, err := session.Model(model).
		Where("id = ? and u_id=?",id,userId).Count()
	if err != nil {
		return false,err
	}
	if count>0{
		return true,nil
	}
	return false,nil
}