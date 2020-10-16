package permission

import "time"

type Permission struct {
	AuthorityID   int       `json:"authorityId"`
	AuthorityName string    `json:"AuthorityName"`
	OrderNumber   int       `json:"orderNumber"`
	MenuURL       string    `json:"menuUrl"`
	MenuIcon      string    `json:"menuIcon"`
	Authority     string    `json:"authority"`
	Checked       int       `json:"checked"`
	IsMenu        int       `json:"isMenu"`
	ParentID      int       `json:"parentId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// 设置Permission的表名为`permission`
func (Permission) TableName() string {
	return "permission"
}
