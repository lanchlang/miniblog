package model

type Category struct {
	tableName struct{} `sql:"t_category"`
	Id          int `sql:"id,pk" json:"id" form:"id"`
	Name        string `sql:"name" json:"name" form:"name"`
	ParentId    int `sql:"p_id" json:"parent_id" form:"parent_id"`
	ParentName  string `sql:"p_name" json:"parent_name" form:"parent_name"`
}
func NewCategory()interface{}{
	return new(Category)
}

func NewCategoryList()interface{}{
	var list []Category
	return &list
}

//获取category的总量
func (model *Category) TotalCount() (int,error) {
	session := GetSession()
	defer session.Close()
	count,err:= session.Model(model).Count()
	if err!=nil{
		return 0,err
	}
	return count,err
}
//通过category搜索
func (model *Category) GetByName(name string) error{
	session := GetSession()
	defer session.Close()
	return session.Model(model).Where("name=?", name).Select()
}

//获取所有的大类
func (model *Category) GetParentCategory() ([]Category,error){
	session := GetSession()
	defer session.Close()
	var categories []Category
	err:=session.Model(&categories).Where("p_id is null").Select()
	return categories,err
}

//获取大类下的子分类
//输入为父类的Id
//输出为子类
func (model *Category) GetSubCategory(parentId int64) ([]Category,error){
	session := GetSession()
	defer session.Close()
	var categories []Category
	err:=session.Model(&categories).Where("p_id =?",parentId).Select()
	return categories,err
}