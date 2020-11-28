package role

import (
	"log"

	"github.com/jinzhu/gorm"
)

type (
	/**
	  面向接口开发：
	  面向接口开发的好处是要对下面的函数进行测试时，不需要依赖一个全局的db连接，只需要调用NewService传一个db连接参数即可测试
	*/
	RoleService interface {
		GetRoleInfo(search Role) (role Role, err error)
		GetRoleList(page, limit int, search Role) (list []Role, count int, err error)
		RolePermissions(roleID int) (role Role, err error)
		CreateRole(role *Role) (err error)
		UpdateRole(ID int, role *Role) (err error)
		DeleteRole(ID int) (err error)
	}
)

type roleService struct {
	db *gorm.DB
}

/*NewService 初始化结构体*/
func NewService(db *gorm.DB) RoleService {
	return &roleService{
		db: db,
	}
}

func (u *roleService) GetRoleInfo(search Role) (role Role, err error) {
	err = u.db.Where(search).Find(&role).Error
	if err != nil {
		return role, err
	}
	return role, nil
}

func (u *roleService) GetRoleList(page, limit int, search Role) (list []Role, count int, err error) {
	err = u.db.Model("role").Order("role_num asc").Find(&list).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return list, count, err
	}
	return list, count, nil
}

//获取角色拥有的权限
func (u *roleService) RolePermissions(roleNum int) (role Role, err error) {
	err = u.db.Model("role").Select("permissions").Where("role_num = ?", roleNum).Order("role_num asc").Find(&role).Error
	if err != nil {
		return role, err
	}
	return role, nil
}

func (u *roleService) CreateRole(role *Role) (err error) {
	err = u.db.Create(role).Error
	log.Println(role)
	if err != nil {
		return err
	}
	return nil
}

func (u *roleService) UpdateRole(ID int, role *Role) (err error) {
	err = u.db.Model(role).Where("id= ? ", ID).Updates(role).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *roleService) DeleteRole(ID int) (err error) {
	err = u.db.Where("id = ?", ID).Delete(Role{}).Error
	if err != nil {
		return err
	}
	return nil
}
