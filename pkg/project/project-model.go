package project

import "time"

type Project struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    int       `json:"status" gorm:"default:1"`
	Env       string    `json:"env"`
	Namespace string    `json:"namespace"`
	GitRepo   string    `json:"gitRepo"`
	Branch    string    `json:"branch"`
	Config    string    `json:"config"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 设置User的表名为`user`
func (Project) TableName() string {
	return "project"
}
