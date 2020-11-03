package permission

type Permission struct {
	ID          int    `json:"id"`          //id
	Title       string `json:"title"`       //权限名称
	OrderNumber int    `json:"orderNumber"` //排序
	MenuURL     string `json:"menuUrl"`     //菜单URL地址
	MenuIcon    string `json:"menuIcon"`    //icon
	Authority   string `json:"authority"`   //权限标识
	Checked     int    `json:"checked"`     //
	IsMenu      int    `json:"isMenu"`      //1按钮，0菜单
	ParentID    int    `json:"parentId"`    //父级id
}

// 设置Permission的表名为`permission`
func (Permission) TableName() string {
	return "permission"
}

type PermissionTree struct {
	Title    string                    `json:"title"`
	ID       int                       `json:"id"`
	Field    string                    `json:"field"`
	Spread   bool                      `json:"spread"`
	Checked  bool                      `json:"checked"`
	Children []*PermissionTreeChildren `json:"children"`
}

type PermissionTreeChildren struct {
	Title   string `json:"title"`
	ID      int    `json:"id"`
	Field   string `json:"field"`
	Checked bool   `json:"checked"`
}
