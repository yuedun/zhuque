package task

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type (
	/**
	  面向接口开发：
	  面向接口开发的好处是要对下面的函数进行测试时，不需要依赖一个全局的mysql连接，只需要调用NewService传一个mysql连接参数即可测试
	*/
	TaskService interface {
		GetTaskInfo(search Task) (task Task, err error)
		GetTaskList(search Task) (list []Task, err error)
		GetTaskInfoBySQL() (task Task, err error)
		CreateTask(task *Task) (err error)
		UpdateTask(serverID int, task *Task) (err error)
		DeleteTask(serverID int) (err error)
	}
)

type projectService struct {
	mysql *gorm.DB
}

/*NewService 初始化结构体*/
func NewService(mysql *gorm.DB) TaskService {
	return &projectService{
		mysql: mysql,
	}
}

func (u *projectService) GetTaskInfo(search Task) (user Task, err error) {
	err = u.mysql.Where(search).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *projectService) GetTaskList(search Task) (list []Task, err error) {
	err = u.mysql.Where(search).Find(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
}

func (u *projectService) GetTaskInfoBySQL() (task Task, err error) {
	err = u.mysql.Raw("select * from task where id=?", task.ID).Scan(&task).Error
	if err != nil {
		return task, err
	}
	return task, nil
}

func (u *projectService) CreateTask(task *Task) (err error) {
	err = u.mysql.Create(task).Error
	fmt.Println(task)
	if err != nil {
		return err
	}
	return nil
}

func (u *projectService) UpdateTask(userID int, task *Task) (err error) {
	err = u.mysql.Model(task).Where("id = ?", userID).Update(task).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *projectService) DeleteTask(ID int) (err error) {
	u.mysql.Where("id = ?", ID).Delete(Task{})
	if err != nil {
		return err
	}
	return nil
}
