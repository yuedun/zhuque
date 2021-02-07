package task

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type (
	/**
	  面向接口开发：
	  面向接口开发的好处是要对下面的函数进行测试时，不需要依赖一个全局的db连接，只需要调用NewService传一个db连接参数即可测试
	*/
	TaskService interface {
		GetTaskInfo(search Task) (task Task, err error)
		GetTaskList(offet, limit int, search Task) (list []Task, count int, err error)
		GetTaskInfoBySQL() (task Task, err error)
		CreateTask(task *Task) (ID int, err error)
		UpdateTask(ID int, task *Task) (err error)
		DeleteTask(ID int) (err error)
		ReleaseTask(ID int) (string, error)
		ReleaseTaskV2(ID int) (string, error)
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

func (u *taskService) GetTaskList(offset, limit int, search Task) (list []Task, count int, err error) {
	err = u.db.Where(search).Offset(offset).Limit(limit).Order("id desc").Find(&list).Offset(-1).Limit(-1).Count(&count).Error //排序
	if err != nil {
		return list, count, err
	}
	return list, count, nil
}

func (u *taskService) GetTaskInfoBySQL() (task Task, err error) {
	err = u.db.Raw("select * from task where id=?", task.ID).Scan(&task).Error
	if err != nil {
		return task, err
	}
	return task, nil
}

func (u *taskService) CreateTask(task *Task) (ID int, err error) {
	err = u.db.Create(task).Error
	log.Println(task)
	if err != nil {
		return 0, err
	}
	return task.ID, nil
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

func (u *taskService) ReleaseTask(ID int) (string, error) {
	var err error
	search := Task{ID: ID}
	task, err := u.GetTaskInfo(search)
	if err != nil {
		return "", err
	}
	var cmdOut []byte
	var cmd *exec.Cmd
	// 执行单个shell命令时, 直接运行即可
	// 从数据库中取出project和cmd组合。
	log.Println("执行命令", task.Cmd)
	cmd = exec.Command("bash", "-c", task.Cmd)
	if cmdOut, err = cmd.CombinedOutput(); err != nil {
		log.Println("输出错误：", err)
		log.Println("输出错误2：", string(cmdOut))
		//保存失败发布记录
		if err = u.db.Model(&task).UpdateColumn("releaseState", 0).Error; err != nil {
			return "更新数据失败", err
		}
		return strings.ReplaceAll(string(cmdOut), "\n", "<br>"), err
	}
	// 默认输出有一个换行
	log.Println(string(cmdOut))
	if err = u.db.Model(&task).UpdateColumn("releaseState", 1).Error; err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(cmdOut), "\n", "<br>"), nil
}

func (u *taskService) ReleaseTaskV2(ID int) (string, error) {
	var err error
	search := Task{ID: ID}
	task, err := u.GetTaskInfo(search)
	if err != nil {
		return "", err
	}
	projectList := strings.Split(task.Project, ",")
	projectLen := len(projectList)
	ch := make(chan string, projectLen)
	for _, projectName := range projectList {
		log.Println("projectName", projectName)
		go excuteCmd(projectName, task.Cmd, ch)
	}
	resultAll := ""
	i := 0
Loop:
	for {
		select {
		case result := <-ch:
			i++
			resultAll += result
			log.Println("projectLen------------------------------", i)
		case <-time.After(time.Second * 300):
			log.Println("timeout!!")
		}
		if projectLen == i {
			break Loop
		}
	}
	log.Println("end for")
	if err = u.db.Model(&task).UpdateColumn("releaseState", 1).Error; err != nil {
		return "更新数据库失败", err
	}
	return resultAll, nil
}

// excuteCmd 用于异步并行执行
// @projectName 项目名 @cmd 执行命令 @result 命令执行结果
func excuteCmd(projectName string, taskCmd string, result chan string) {
	var cmdOut []byte
	var cmd *exec.Cmd
	var err error
	taskCmd = fmt.Sprintf(taskCmd, projectName)
	log.Println("执行命令", taskCmd)
	cmd = exec.Command("bash", "-c", taskCmd)
	if cmdOut, err = cmd.CombinedOutput(); err != nil {
		log.Println(projectName+"输出错误：", err)
		log.Println(projectName+"输出错误2：", string(cmdOut))
	}
	result <- strings.ReplaceAll(string(cmdOut), "\n", "<br>") //写入命令执行结果
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
