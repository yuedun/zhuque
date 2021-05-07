package role

type Role struct {
	ID          int    `json:"id" gorm:"AUTO_INCREMENT;primary_key"`
	RoleNum     int    `json:"roleNum" gorm:"comment:'超管1，管理员2，开发3'"` //超管1，管理员2，开发3
	Name        string `json:"name"`                                  //说明
	Permissions string `json:"permissions"`                           //该角色拥有的权限
}

// 设置Permission的表名为`permission`
func (Role) TableName() string {
	return "role"
}
