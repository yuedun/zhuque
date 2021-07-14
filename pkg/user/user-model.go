package user

import "time"

type User struct {
	ID        int       `json:"id" gorm:"AUTO_INCREMENT;primary_key"`
	UserName  string    `json:"username" gorm:"type:varchar(10)"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Status    int       `json:"status"`
	RoleNum   int       `json:"roleNum" grom:"comment:'1超级管理员（最高权限），2管理员（管理空间，一个空间有多个项目），3发布人员'"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

const (
	SuperManager = 1
	Manager      = 2
	Developer    = 3
)

// 设置User的表名为`user`
func (User) TableName() string {
	return "user"
}

type UserLoginInfo struct {
	ID          int    `json:"id"`
	UserName    string `json:"username"`
	Password    string `json:"password"`
	RoleNum     int    `json:"roleNum"` // 1超级管理员（最高权限），2管理员（管理空间，一个空间有多个项目），3发布人员
	Permissions string `json:"permissions"`
}
