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
		GetTask(search Task) (task Task, err error)
		GetTaskInfo(search Task) (task Task, err error)
		GetTaskList(offet, limit int, search Task) (list []Task, count int, err error)
		WaitTaskList(from string) (list []Task, err error)
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

// 查询任务，关联查出用户
func (u *taskService) GetTask(search Task) (task Task, err error) {
	err = u.db.Where(search).First(&task).Error
	err = u.db.Model(&task).Related(&task.User).Error
	if err != nil {
		return task, err
	}
	return task, nil
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

func (u *taskService) WaitTaskList(from string) (list []Task, err error) {
	err = u.db.Raw("select * from task where `from` = ? and release_state in (?) order by id desc;", from, []int{2, 3}).Scan(&list).Error //排序
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
	userCmd := fmt.Sprintf("pm2 deploy projects/%s/ecosystem.config.js production --force", task.Project)
	log.Println("执行命令", userCmd)
	if err = u.db.Model(&task).UpdateColumn("releaseState", Releaseing).Error; err != nil {
		return "更新数据失败", err
	}
	cmd = exec.Command("bash", "-c", userCmd)
	if cmdOut, err = cmd.CombinedOutput(); err != nil {
		log.Println("输出错误：", err)
		log.Println("输出错误2：", string(cmdOut))
		//保存失败发布记录
		if err = u.db.Model(&task).UpdateColumn("releaseState", Fail).Error; err != nil {
			return "更新数据失败", err
		}
		return strings.ReplaceAll(string(cmdOut), "\n", "<br>"), err
	}
	// 默认输出有一个换行
	log.Println(">>>>>", string(cmdOut))
	if err = u.db.Model(&task).UpdateColumn("releaseState", Success).Error; err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(cmdOut), "\n", "<br>"), nil
}

type CmdResult struct {
	Status  int
	Content string
}

func (u *taskService) ReleaseTaskV2(ID int) (string, error) {
	var err error
	search := Task{ID: ID}
	task, err := u.GetTaskInfo(search)
	if err != nil {
		return "", err
	}
	if err = u.db.Model(&task).UpdateColumn("releaseState", Releaseing).Error; err != nil {
		return "更新数据失败", err
	}
	projectList := strings.Split(task.Project, ",")
	projectLen := len(projectList)
	ch := make(chan CmdResult, projectLen)
	for i, projectName := range projectList {
		userCmd := fmt.Sprintf("pm2 deploy projects/%s/ecosystem.config.js production --force", projectName)
		log.Println("执行命令", i, userCmd)
		go excuteCmd(userCmd, ch)
	}
	resultAll := ""
	i := 0
	result := CmdResult{}
Loop:
	for {
		select {
		case result = <-ch:
			i++
			resultAll += result.Content
			log.Println("projectLen------------------------------", i)
		case <-time.After(time.Second * 300):
			log.Println("timeout!!")
		}
		if projectLen == i {
			break Loop
		}
	}
	log.Println("end for")
	if err = u.db.Model(&task).UpdateColumn("releaseState", result.Status).Error; err != nil {
		return "更新数据库失败", err
	}
	return resultAll, nil
}

// excuteCmd 用于异步并行执行
// @taskCmd 执行命令 @result 命令执行结果
func excuteCmd(taskCmd string, result chan CmdResult) {
	cmdResult := CmdResult{
		Status:  Success,
		Content: "",
	}
	var cmdOut []byte
	var cmd *exec.Cmd
	var err error
	log.Println("输入命令：", taskCmd)
	cmd = exec.Command("bash", "-c", taskCmd)
	if cmdOut, err = cmd.CombinedOutput(); err != nil {
		log.Println("输出错误：", err)
		log.Println("输出错误2：", string(cmdOut))
		cmdResult.Status = Fail
	}
	cmdResult.Content = strings.ReplaceAll(string(cmdOut), "\n", "<br>") //写入命令执行结果
	result <- cmdResult
}

// .NowRelease
func (u *taskService) Approve(params map[string]interface{}) (task Task, err error) {
	id64 := params["id"].(float64)
	id := int(id64)
	task.ID = id
	err = u.db.Model(&task).UpdateColumns(params).Find(&task).Error
	if err != nil {
		return task, err
	}
	return task, nil
}
