package task

import "time"

type Task struct {
	ID        int       `json:"id"`
	TaskName  string    `json:"taskName"`
	Project   string    `json:"project"`
	Username  string    `json:"username"`
	CommitID  string    `json:"commitID"`
	Env       string    `json:"env"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 设置User的表名为`user`
func (Task) TableName() string {
	return "task"
}
