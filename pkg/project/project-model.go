package project

import "time"

type Project struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Env       string    `json:"env"`
	Namespace string    `json:"namespace"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 设置User的表名为`user`
func (Project) TableName() string {
	return "project"
}
