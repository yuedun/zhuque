package user

type UserProject struct {
	ID         int `json:"id"`
	UserID     int `json:"userID"`
	ProjectID  int `json:"projectID"`
	CreateUser int `json:"createUser"`
}

// 设置UserProject的表名为`user_project`
func (UserProject) TableName() string {
	return "user_project"
}

type UserProjectVO struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`       //项目名
	Namespace  string `json:"namespace"`  //空间
	Username   string `json:"username"`   //关联用户
	CreateUser string `json:"createUser"` //创建者
}
