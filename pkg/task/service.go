package task

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/jinzhu/gorm"
)

type (
	/**
	  面向接口开发：
	  面向接口开发的好处是要对下面的函数进行测试时，不需要依赖一个全局的db连接，只需要调用NewService传一个db连接参数即可测试
	*/
	TaskService interface {
		GetTaskInfo(search Task) (task Task, err error)
		GetTaskList(search Task) (list []Task, err error)
		GetTaskInfoBySQL() (task Task, err error)
		CreateTask(task *Task) (err error)
		UpdateTask(ID int, task *Task) (err error)
		DeleteTask(ID int) (err error)
		ReleaseTask(ID int) (error, string)
		Approve(params map[string]interface{}) (Task, error)
	}
)

type taskService struct {
	db *gorm.DB
}

/*NewService 初始化结构体*/
func NewService(db *gorm.DB) TaskService {
	return &taskService{
		db: db,
	}
}

func (u *taskService) GetTaskInfo(search Task) (task Task, err error) {
	err = u.db.Where(search).Find(&task).Error
	if err != nil {
		return task, err
	}
	return task, nil
}

func (u *taskService) GetTaskList(search Task) (list []Task, err error) {
	err = u.db.Where(search).Order("id desc").Find(&list).Error //排序
	if err != nil {
		return list, err
	}
	return list, nil
}

func (u *taskService) GetTaskInfoBySQL() (task Task, err error) {
	err = u.db.Raw("select * from task where id=?", task.ID).Scan(&task).Error
	if err != nil {
		return task, err
	}
	return task, nil
}

func (u *taskService) CreateTask(task *Task) (err error) {
	err = u.db.Create(task).Error
	fmt.Println(task)
	if err != nil {
		return err
	}
	return nil
}

func (u *taskService) UpdateTask(ID int, task *Task) (err error) {
	err = u.db.Model(task).UpdateColumn("releaseState", task.ReleaseState).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *taskService) DeleteTask(ID int) (err error) {
	err = u.db.Where("id = ?", ID).Delete(Task{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *taskService) ReleaseTask(ID int) (error, string) {
	var err error
	search := Task{ID: ID}
	task, err := u.GetTaskInfo(search)
	if err != nil {
		return err, ""
	}
	var cmdOut []byte
	var cmd *exec.Cmd
	// 执行单个shell命令时, 直接运行即可
	cmd = exec.Command("bash", "-c", task.Cmd)
	if cmdOut, err = cmd.CombinedOutput(); err != nil {
		log.Println("输出错误：", err)
		log.Println("输出错误2：", string(cmdOut))
		//保存失败发布记录
		if err = u.db.Model(&task).UpdateColumn("releaseState", 0).Error; err != nil {
			return err, "更新数据失败"
		}
		return err, strings.ReplaceAll(string(cmdOut), "\n", "<br>")
	}
	// 默认输出有一个换行
	log.Println(string(cmdOut))
	if err = u.db.Model(&task).UpdateColumn("releaseState", 1).Error; err != nil {
		return err, ""
	}
	return nil, strings.ReplaceAll(string(cmdOut), "\n", "<br>")
}

func (u *taskService) Approve(params map[string]interface{}) (task Task, err error) {
	id64 := params["id"].(float64)
	id := int(id64)
	task.ID = id
	err = u.db.Model(&Task{ID: id}).UpdateColumns(params).Find(&task).Error
	if err != nil {
		return task, err
	}
	return task, nil
}
