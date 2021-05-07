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
	User       string   `json:"user"`
	Host       []string `json:"host"`
	Ref        string   `json:"ref"`
	Repo       string   `json:"repo"`
	Path       string   `json:"path"`
	PreBuild   string   `json:"preBuild"`   //构建前执行的命令：装依赖，设置环境变量等
	Build      string   `json:"build"`      // 构建命令
	PreDeploy  string   `json:"preDeploy"`  //应用服务器重启前执行的命令：设置环境变量等
	PostDeploy string   `json:"postDeploy"` //应用服务器重启
	RsyncArgs  string   `json:"rsyncArgs"`  // rsync参数
}

// 设置User的表名为`user`
func (Project) TableName() string {
	return "project"
}
