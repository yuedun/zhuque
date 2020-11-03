package permission

import (
	"fmt"
	"log"
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
		GetPermissionList(userId int) (list []Permission, count int, err error)
		CreatePermission(permission *Permission) (err error)
		UpdatePermission(ID int, permission *Permission) (err error)
		DeletePermission(ID int) (err error)
		GetByRole(roleID int) (list []Permission, err error)
		GetPermissionForSide(roleNum int) (list []Permission, err error)
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

func (u *permissionService) GetPermissionList(userID int) (list []Permission, count int, err error) {
	err = u.db.Model("permission").Order("order_number asc").Find(&list).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return list, count, err
	}
	return list, count, nil
}

// GetPermissionForSide 侧边栏菜单
func (u *permissionService) GetPermissionForSide(roleNum int) (menus []Permission, err error) {
	// roleService := role.NewService(db.SQLLite)
	// roleObj := role.Role{
	// 	RoleNum: roleNum,
	// }
	// roleResult, _ := roleService.GetRoleInfo(roleObj)
	//查询所有父级菜单
	err = u.db.Model("permission").Where("is_menu = 0 AND parent_id > 0").Find(&menus).Error
	if err != nil {
		return menus, err
	}
	return menus, nil
}

func (u *permissionService) CreatePermission(permission *Permission) (err error) {
	err = u.db.Create(permission).Error
	fmt.Println(permission)
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
	log.Println(">>>>>>>>>>permissions", permissions)
	for _, pID := range permissions {
		permis := Permission{}
		log.Println(">>>>>>>>>>pid", pID)
		err = u.db.Model("permission").Select("id, authority_name").Where("id = ?", pID).Find(&permis).Error
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
