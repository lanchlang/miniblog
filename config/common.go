package config

type CommonConfig struct{
	DefaultListSize int `json:"default_list_size"`  //默认列表长度
	JwtToken        string `json:"jwt_token"`     //jwt token的名字
	User        string `json:"user"`     //jwt token的名字
	BlogAdminAccessLevel int `json:"blog_admin_access_level"` //blog需要的管理级别
	CommentAdminAccessLevel int `json:"comment_admin_access_level"` //comment需要的管理级别
	UserAdminAccessLevel int `json:"user_admin_access_level"` //用户管理员需要的管理级别
	DefaultUserAccessLevel int `json:"default_user_access_level"` //注册用户需要的管理级别
}
var DefaultConfig=new(CommonConfig)
func init() {
    err:=LoadConfig("config/common.json",DefaultConfig)
    if err!=nil{
    	panic(err)
	}
}
