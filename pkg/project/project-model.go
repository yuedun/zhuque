package project

import "time"

type Project struct {
	ID   int    `json:"id" gorm:"AUTO_INCREMENT;primary_key"`
	Name string `json:"name" gorm:"type:varchar(20);unique_index"`
	// Clone      string    `json:"clone" gorm:"type:varchar(20);comment:'同项目不同环境'"`
	Status     int       `json:"status" gorm:"default:1;not null;comment:'1正常，0删除'"`
	Env        string    `json:"env"`
	Namespace  string    `json:"namespace"`
	Config     string    `json:"config" gorm:"comment:'项目配置'"`                //发布配置
	DeployType string    `json:"deployType" grom:"comment:'发布机制，值：pm2 ,scp'"` // 发布机制，值：pm2 ,scp
	DingTalk   string    `json:"dingTalk" gorm:"comment:'钉钉hook地址'"`
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
