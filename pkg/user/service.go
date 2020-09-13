package user

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type (
	/**
	  面向接口开发：
	  面向接口开发的好处是要对下面的函数进行测试时，不需要依赖一个全局的mysql连接，只需要调用NewService传一个mysql连接参数即可测试
	*/
	UserService interface {
		GetUserInfo(userObj User) (user User, err error)
		GetUserList(offset, limit int, userObj User) (user []User, count int, err error)
		GetUserProjects(offset, limit int, userObj User) (userProjects []UserProjectVO, count int, err error)
		GetUserInfoBySQL() (user User, err error)
		CreateUser(user *User) (err error)
		CreateUserProject(userProject *UserProject) (err error)
		UpdateUser(userID int, user *User) (err error)
		DeleteUser(userID int) (err error)
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

func (u *userService) GetUserList(offset, limit int, userObj User) (users []User, count int, err error) {
	err = u.mysql.Where(userObj).Offset(offset).Limit(limit).Find(&users).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return users, count, err
	}
	return users, count, nil
}

func (u *userService) GetUserProjects(offset, limit int, userObj User) (userProjects []UserProjectVO, count int, err error) {
	err = u.mysql.Table("user_project AS up").
		Select("up.id AS id, p.name AS name, p.namespace AS namespace, u.user_name AS username, u2.user_name AS createUser").
		Joins("INNER JOIN user AS u ON u.id = up.user_id").
		Joins("INNER JOIN project AS p ON p.id = up.project_id").
		Joins("INNER JOIN user AS u2 ON up.create_user = u2.id").
		Where("u.id = ?", userObj.ID).Find(&userProjects).Error
	if err != nil {
		return userProjects, count, err
	}
	return userProjects, count, nil
}

func (u *userService) GetUserInfoBySQL() (user User, err error) {
	err = u.mysql.Raw("select * from user where id=?", user.ID).Scan(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *userService) CreateUser(user *User) (err error) {
	err = u.mysql.Create(user).Error
	fmt.Println(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) CreateUserProject(userProject *UserProject) (err error) {
	err = u.mysql.Create(userProject).Error
	fmt.Println(userProject)
	if err != nil {
		return err
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
