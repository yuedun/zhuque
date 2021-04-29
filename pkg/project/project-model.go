package project

import "time"

type Project struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Status     int       `json:"status" gorm:"default:1"`
	Env        string    `json:"env"`
	Namespace  string    `json:"namespace"`
	Config     string    `json:"config"`     // pm2发布配置
	DeployType string    `json:"deployType"` // 发布机制，值：pm2 ,scp
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type DeployConfig struct {
	APPName        string   `json:"appName"` // 用于clone代码指定项目名
	User           string   `json:"user"`
	Host           []string `json:"host"`
	Ref            string   `json:"ref"`
	Repo           string   `json:"repo"`
	Path           string   `json:"path"`
	PreDeployLocal string   `json:"pre-deploy-local"`
	PostDeploy     string   `json:"post-deploy"`
	PreSetup       string   `json:"pre-setup"`
	Build          string   `json:"build"`
	RsyncArgs      string   `json:"rsyncArgs"` // rsync参数
}

// 设置User的表名为`user`
func (Project) TableName() string {
	return "project"
}
