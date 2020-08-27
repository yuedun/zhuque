package project

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type (
	/**
	  面向接口开发：
	  面向接口开发的好处是要对下面的函数进行测试时，不需要依赖一个全局的mysql连接，只需要调用NewService传一个mysql连接参数即可测试
	*/
	ProjectService interface {
		GetProjectInfo(search Project) (project Project, err error)
		GetProjectList(offset, limit int, search Project) (list []Project, count int, err error)
		GetProjectInfoBySQL() (project Project, err error)
		CreateProject(project *Project) (err error)
		UpdateProject(serverID int, project *Project) (err error)
		DeleteProject(serverID int) (err error)
	}
)

type projectService struct {
	mysql *gorm.DB
}

/*NewService 初始化结构体*/
func NewService(mysql *gorm.DB) ProjectService {
	return &projectService{
		mysql: mysql,
	}
}

func (u *projectService) GetProjectInfo(search Project) (user Project, err error) {
	err = u.mysql.Where(search).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *projectService) GetProjectList(offset, limit int, search Project) (list []Project, count int, err error) {
	err = u.mysql.Where(search).Offset(offset).Limit(limit).Find(&list).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return list, count, err
	}
	return list, count, nil
}

func (u *projectService) GetProjectInfoBySQL() (project Project, err error) {
	err = u.mysql.Raw("select * from project where id=?", project.ID).Scan(&project).Error
	if err != nil {
		return project, err
	}
	return project, nil
}

func (u *projectService) CreateProject(project *Project) (err error) {
	err = u.mysql.Create(project).Error
	fmt.Println(project)
	if err != nil {
		return err
	}
	return nil
}

func (u *projectService) UpdateProject(userID int, project *Project) (err error) {
	err = u.mysql.Model(project).Where("id = ?", userID).Update(project).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *projectService) DeleteProject(ID int) (err error) {
	u.mysql.Where("id = ?", ID).Delete(Project{})
	if err != nil {
		return err
	}
	return nil
}
