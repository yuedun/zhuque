package permission

type Permission struct {
	ID          int    `json:"id"`          //id
	Title       string `json:"title"`       //权限名称
	OrderNumber int    `json:"orderNumber"` //排序，父菜单为1，则子菜单按照：11,12,13排序，父菜单为2，则子菜单按照：21,22,23排序。两位数够了。
	Href        string `json:"href"`        //菜单URL地址
	Icon        string `json:"icon"`        //icon
	Authority   string `json:"authority"`   //权限标识
	Checked     int    `json:"checked"`     //
	MenuType    int    `json:"menuType"`    //1按钮，0菜单
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

// 顶部和侧边菜单
type Menus struct {
	Permission
	Target string       `json:"target"` //没啥用
	Child  []Permission `json:"child"`
}
