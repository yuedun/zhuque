package task

import "time"

type Task struct {
	ID           int       `json:"id"`
	TaskName     string    `json:"taskName"`
	Project      string    `json:"project"`
	UserID       string    `json:"userID"`
	Username     string    `json:"username"`
	CommitID     string    `json:"commitID"`
	Env          string    `json:"env"`
	Status       string    `json:"status"`       //数据状态：1有效，0无效
	ReleaseState int       `json:"releaseState"` // 发布结果：1成功，0失败，2待发布
	Cmd          string    `json:"cmd"`          // 发布命令
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// 设置User的表名为`user`
func (Task) TableName() string {
	return "task"
}
