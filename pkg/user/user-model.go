package user

import "time"

type User struct {
	Id          int       `json:"id"`
	Mobile      string    `json:"mobile"`
	UserName    string    `json:"userName"`
	Password    string    `json:"password"`
	Gender      string    `json:"gender"`
	Addr        string    `json:"addr"`
	Email       string    `json:"email"`
	Status      int       `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// 设置User的表名为`user`
func (User) TableName() string {
	return "user"
}
