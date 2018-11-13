package model

import (
	"time"
	"github.com/go-pg/pg/orm"
)

type Verification struct {
	Id      int64
	Email   string
	Phone   string
	Code    string
	Timeout time.Time
}

func (model *Verification) BeforeInsert(db orm.DB) error {
	if len(model.Email) > 0 && model.Timeout.IsZero() {
		model.Timeout = time.Now().Add(72 * time.Hour)
	}
	return nil
}
func (model *Verification) Save() error {
	session := GetSession()
	defer session.Close()
	return session.Insert(model)
}
func (model *Verification) Delete(id int64) error {
	session := GetSession()
	defer session.Close()
	_, err := session.Model(model).Where("id=?", id).Delete()
	if err != nil {
		return err
	} else {
		return nil
	}
}
func (model *Verification) Find(code string) error {
	session := GetSession()
	defer session.Close()
	return session.Model(model).Where("code=?", code).Select()
}

