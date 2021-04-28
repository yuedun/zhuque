package role

type RolePermission struct {
	ID           int    `json:"id"`
	RoleID       int    `json:"roleID"`
	PermissionID string `json:"permissionID"`
}

// 设置Permission的表名为`permission`
func (RolePermission) TableName() string {
	return "role_permission"
}
