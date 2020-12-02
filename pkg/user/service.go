package user

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
)

type (
	/**
	  面向接口开发：
	  面向接口开发的好处是要对下面的函数进行测试时，不需要依赖一个全局的mysql连接，只需要调用NewService传一个mysql连接参数即可测试
	*/
	UserService interface {
		GetUserInfo(userObj User) (user User, err error)
		//获取登录用户信息，包含权限信息
		GetLoginUserInfo(userObj User) (user UserLoginInfo, err error)
		GetUserList(offset, limit int, userObj User) (user []User, count int, err error)
		GetUserProjects(offset, limit int, userObj User) (userProjects []UserProjectVO, count int, err error)
		CreateUser(user *User) (err error)
		CreateUserProject(userProject *UserProject) (err error)
		UpdateUser(userID int, user *User) (err error)
		DeleteUser(userID int) (err error)
		DeleteUserProject(upID int) (err error)
		//获取项目关联用户
		GetProjectUsersEmail(projectName string) (emails []string, err error)
	}
)

type userService struct {
	mysql *gorm.DB
}

/*NewService 初始化结构体*/
func NewService(mysql *gorm.DB) UserService {
	return &userService{
		mysql: mysql,
	}
}

func (u *userService) GetUserInfo(userObj User) (user User, err error) {
	err = u.mysql.Where(userObj).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
func (u *userService) GetLoginUserInfo(userObj User) (user UserLoginInfo, err error) {
	err = u.mysql.Table("user AS u").
		Select("u.id AS id, u.user_name AS user_name, u.password AS password, r.permissions AS permissions").
		Joins("INNER JOIN role AS r ON u.role_num = r.role_num").
		Where("u.user_name = ?", userObj.UserName).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *userService) GetUserList(offset, limit int, userObj User) (users []User, count int, err error) {
	err = u.mysql.Where(userObj).Offset(offset).Limit(limit).Find(&users).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return users, count, err
	}
	return users, count, nil
}

func (u *userService) GetUserProjects(offset, limit int, userObj User) (userProjects []UserProjectVO, count int, err error) {
	// 字段别名需要设置成下划线命名法，不能设置为驼峰
	err = u.mysql.Table("user_project AS up").
		Select("up.id AS id, p.name AS name, p.namespace AS namespace, u.user_name AS username, u2.user_name AS create_user").
		Joins("INNER JOIN user AS u ON u.id = up.user_id").
		Joins("INNER JOIN project AS p ON p.id = up.project_id").
		Joins("INNER JOIN user AS u2 ON up.create_user = u2.id").
		Where("u.id = ?", userObj.ID).Find(&userProjects).Error
	if err != nil {
		return userProjects, count, err
	}
	return userProjects, count, nil
}

func (u *userService) CreateUser(user *User) (err error) {
	err = u.mysql.Create(user).Error
	log.Println(user)
	if err != nil {
		return err
	}
	return nil
}

// CreateUserProject 先查询是否存在，再创建
func (u *userService) CreateUserProject(search *UserProject) (err error) {
	err = u.mysql.Model("user_project").Where("user_id=? and project_id = ?", search.UserID, search.ProjectID).Find(search).Error
	if err != nil {
		log.Println(">>>>>>>>>无数据", err)
		err = u.mysql.Create(search).Error
		if err != nil {
			return err
		}
	} else {
		return errors.New("用户已存在该项目！")
	}
	return nil
}

func (u *userService) UpdateUser(userID int, user *User) (err error) {
	err = u.mysql.Model(user).Where("id = ?", userID).Update(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) DeleteUser(userID int) (err error) {
	u.mysql.Where("id = ?", userID).Delete(User{})
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) DeleteUserProject(upID int) (err error) {
	u.mysql.Where("id = ?", upID).Delete(UserProject{})
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) GetProjectUsersEmail(projectName string) (emails []string, err error) {
	// SELECT u.email from project AS p INNER JOIN user_project AS up ON p.id = up.project_id INNER JOIN user AS u ON u.id = up.user_id WHERE p.name='zhuque'
	err = u.mysql.Table("project AS p").
		Select("u.email AS email").
		Joins("INNER JOIN user_project AS up ON p.id = up.project_id").
		Joins("INNER JOIN user AS u ON u.id = up.user_id").
		Where("p.name = ?", projectName).Pluck("email", &emails).Error
	if err != nil {
		return emails, err
	}
	return emails, nil
}
