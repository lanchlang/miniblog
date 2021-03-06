package model

import (
	"github.com/go-pg/pg/orm"
	"time"
)

type Comment struct {
	tableName struct{} `sql:"t_comment"`
	Id           int64 `sql:"id,pk" json:"id" form:"id"`
	UserId       int64 `sql:"u_id" json:"user_id" form:"user_id"`
	Username     string `sql:"u_name" json:"username" form:"username"`
	BlogId       int64  `sql:"b_id" json:"blog_id" form:"blog_id"`
	Content      map[string]interface{} `sql:"content" json:"content" form:"content"`
	CreateDate   time.Time  `sql:"date" json:"create_date" form:"create_date"`
	Type         int8  `sql:"type" json:"type" form:"type"`  //类型
	ReplyId       int64 `sql:"reply_u_id" json:"reply_id" form:"reply_id"`
	ReplyUsername     string `sql:"reply_u_name" json:"reply_username" form:"reply_username"`
	Deleted       int8  `sql:"deleted" json:"deleted" form:"deleted"`   //deleted大于0，则表示已经删除
}
func NewComment()interface{}{
	return new(Comment)
}

func NewCommentList()interface{}{
	var list []Comment
	return &list
}
func (comment *Comment) BeforeInsert(db orm.DB) error {
	if comment.CreateDate.IsZero() {
		comment.CreateDate = time.Now()
	}
	return nil
}

//用户删除自己的comment,deleted
//输入为comment的id和用户id
//输出为error
func (model *Comment) Delete(id int64, userId int64) error {
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Set("deleted=1").Where("id=? and u_id=?", id, userId).Update()
	return err
}
//获取blog的评论
//输入为blog的id,上一条评论的id，获取的最大数量limit
//输出为comment数组和error
//返回的comment数组是顺序排列的，最新的comment最后返回
func (model *Comment) GetBlogComment(blogId int64,lastCommentId int64, num int) ([]Comment, error) {
	session := GetSession()
	defer session.Close()
	var comments []Comment
	err := session.Model(&comments).Where("b_id=? and id>?", blogId,lastCommentId).Order("id asc").Limit(num).Select()
	if err != nil {
		return nil, err
	}
	return comments, nil
}

