package model

import (
	"github.com/go-pg/pg/orm"
	"time"
)

type Comment struct {
	tableName struct{} `sql:"t_comment"`
	Id           int64 `sql:"id,pk" json:"id"`
	UserId       int64 `sql:"u_id" json:"user_id"`
	Username     string `sql:"u_name" json:"username"`
	Avatar       map[string]interface{} `sql:"u_avatar" json:"avatar"`
	BlogId       int64  `sql:"b_id" json:"blog_id"`
	Content      map[string]interface{} `sql:"content" json:"content"`
	CreateDate   time.Time  `sql:"date" json:"create_date"`
	Type         int8  `sql:"type" json:"type"`  //类型
	ReplyId       int64 `sql:"reply_u_id" json:"reply_id"`
	ReplyUsername     string `sql:"reply_u_name" json:"reply_username"`
	ReplyAvatar       map[string]interface{} `sql:"reply_u_avatar" json:"reply_avatar"`
	Deleted       int8  `sql:"deleted" json:"deleted"`   //deleted大于0，则表示已经删除
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
	err := session.Model(&comments).Where("b_id=? and id>?", blogId,lastCommentId).Limit(num).Select()
	if err != nil {
		return nil, err
	}
	return comments, nil
}

