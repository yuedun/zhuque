package user

type UserProject struct {
	ID        int `json:"id"`
	UserID    int `json:"userID"`
	ProjectID int `json:"projectID"`
}

// 设置UserProject的表名为`user_project`
func (UserProject) TableName() string {
	return "user_project"
}
