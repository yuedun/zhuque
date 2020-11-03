package permission

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/yuedun/zhuque/pkg/role"
)

type (
	/**
	  面向接口开发：
	  面向接口开发的好处是要对下面的函数进行测试时，不需要依赖一个全局的db连接，只需要调用NewService传一个db连接参数即可测试
	*/
	PermissionService interface {
		GetPermissionInfo(search Permission) (permission Permission, err error)
		GetPermissionList() (list []Permission, err error)
		GetPermissionListForRole() (list []*PermissionTree, err error)
		CreatePermission(permission *Permission) (err error)
		UpdatePermission(ID int, permission *Permission) (err error)
		DeletePermission(ID int) (err error)
		GetByRole(roleID int) (list []Permission, err error)
	}
)

type permissionService struct {
	db *gorm.DB
}

/*NewService 初始化结构体*/
func NewService(db *gorm.DB) PermissionService {
	return &permissionService{
		db: db,
	}
}

func (u *permissionService) GetPermissionInfo(search Permission) (permission Permission, err error) {
	err = u.db.Where(search).Find(&permission).Error
	if err != nil {
		return permission, err
	}
	return permission, nil
}

func (u *permissionService) GetPermissionList() (list []Permission, err error) {
	err = u.db.Model("permission").Order("order_number asc").Find(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
}

func (u *permissionService) GetPermissionListForRole() (tree []*PermissionTree, err error) {
	list := []Permission{}
	//查询所有父级菜单
	err = u.db.Model("permission").Where("is_menu = 0 AND parent_id > 0").Find(&list).Error
	for _, pMenu := range list {
		permis := new(PermissionTree)
		permis.ID = pMenu.ID
		permis.Title = pMenu.Title
		permis.Field = pMenu.Authority
		permisChildrenList := []Permission{}
		//获取子菜单
		err = u.db.Model("permission").Select("id, title").Where("is_menu = 1 AND parent_id = ?", pMenu.ID).Find(&permisChildrenList).Error
		if err != nil {
			return nil, err
		}
		treeChildrenList := []*PermissionTreeChildren{}
		for _, children := range permisChildrenList {
			treeChildre := PermissionTreeChildren{}
			treeChildre.ID = children.ID
			treeChildre.Title = children.Title
			treeChildre.Field = children.Authority
			treeChildrenList = append(treeChildrenList, &treeChildre)
		}
		permis.Children = treeChildrenList
		tree = append(tree, permis)
	}
	if err != nil {
		return tree, err
	}
	return tree, nil
}

func (u *permissionService) CreatePermission(permission *Permission) (err error) {
	err = u.db.Create(permission).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *permissionService) UpdatePermission(ID int, permission *Permission) (err error) {
	err = u.db.Model(permission).Where("authority_id = ?", ID).Updates(permission).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *permissionService) DeletePermission(ID int) (err error) {
	err = u.db.Where("id = ?", ID).Delete(Permission{}).Error
	if err != nil {
		return err
	}
	return nil
}

// 获取角色拥有的权限
func (u *permissionService) GetByRole(roleID int) (list []Permission, err error) {
	role := new(role.Role)
	err = u.db.Model("role").Where("id = ?", roleID).Find(role).Error
	permissions := strings.Split(role.Permissions, ",")
	for _, pID := range permissions {
		permis := Permission{}
		err = u.db.Model("permission").Select("id, title").Where("id = ?", pID).Find(&permis).Error
		if err != nil {
			return nil, err
		}
		list = append(list, permis)
	}
	if err != nil {
		return list, err
	}
	return list, nil
}
