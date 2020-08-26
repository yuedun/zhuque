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
	Status       int       `json:"status"`                          //数据状态：1有效，0无效
	ReleaseState int       `json:"releaseState"`                    // 发布结果：1成功，0失败，2待发布
	NowRelease   bool      `json:"nowRelease" gorm:"default:false"` // 是否可以立即发布，需要等待n分钟后发布，该值由管理员审批操作
	Cmd          string    `json:"cmd"`                             // 发布命令
	ApproveMsg   string    `json:"approveMsg"`                      //审批意见
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// 设置Task的表名为`task`
func (Task) TableName() string {
	return "task"
}
