package project

import (
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
		GetProjectNameList(userID int) (list []Project, err error)
		GetAllProjectNameList() (list []Project, err error)
		GetProjectInfoBySQL() (project Project, err error)
		CreateProject(project *Project) (err error)
		UpdateProject(projectID int, project *Project) (err error)
		DeleteProject(prjectID int) (err error)
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
	if search.Name != "" {
		u.mysql = u.mysql.Where("name LIKE ?", search.Name+"%")
	}
	err = u.mysql.Offset(offset).Limit(limit).Find(&list).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return list, count, err
	}
	return list, count, nil
}

// GetProjectNameList 查询登录用户可以发布的项目。只查询项目空间，项目名字段
func (u *projectService) GetProjectNameList(userID int) (list []Project, err error) {
	err = u.mysql.Table("user_project AS up").
		Select("p.id, p.name, p.namespace, p.deploy_type").
		Joins("INNER JOIN project AS p ON p.id = up.project_id").
		Where("up.user_id = ?", userID).Find(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
}

// GetAllProjectNameList 查询登录用户可以发布的项目。只查询项目空间，项目名字段
func (u *projectService) GetAllProjectNameList() (list []Project, err error) {
	err = u.mysql.Find(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
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
	if err != nil {
		return err
	}
	return nil
}

func (u *projectService) UpdateProject(projectID int, project *Project) (err error) {
	err = u.mysql.Model(project).Where("id = ?", projectID).Update(project).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *projectService) DeleteProject(projectID int) (err error) {
	u.mysql.Model(&Project{ID: projectID}).UpdateColumn("status", 0)
	if err != nil {
		return err
	}
	return nil
}
