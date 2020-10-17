package permission

import "time"

type Permission struct {
	ID            int       `json:"id"`            //id
	AuthorityName string    `json:"authorityName"` //权限名称
	OrderNumber   int       `json:"orderNumber"`   //排序
	MenuURL       string    `json:"menuUrl"`       //菜单URL地址
	MenuIcon      string    `json:"menuIcon"`      //icon
	Authority     string    `json:"authority"`     //权限标识
	Checked       int       `json:"checked"`       //
	IsMenu        int       `json:"isMenu"`        //1按钮，0菜单
	ParentID      int       `json:"parentId"`      //父级id
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// 设置Permission的表名为`permission`
func (Permission) TableName() string {
	return "permission"
}
