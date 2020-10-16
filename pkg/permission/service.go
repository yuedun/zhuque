package permission

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type (
	/**
	  面向接口开发：
	  面向接口开发的好处是要对下面的函数进行测试时，不需要依赖一个全局的db连接，只需要调用NewService传一个db连接参数即可测试
	*/
	PermissionService interface {
		GetPermissionInfo(search Permission) (permission Permission, err error)
		GetPermissionList(userId int) (list []Permission, err error)
		CreatePermission(permission *Permission) (err error)
		UpdatePermission(ID int, permission *Permission) (err error)
		DeletePermission(ID int) (err error)
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

func (u *permissionService) GetPermissionList(userId int) (list []Permission, err error) {
	err = u.db.Model("permission").Order("order_number asc").Find(&list).Offset(-1).Limit(-1).Error
	if err != nil {
		return list, err
	}
	return list, nil
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
	err = u.db.Model(permission).UpdateColumn("releaseState", permission.AuthorityID).Error
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
