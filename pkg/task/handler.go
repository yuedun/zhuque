package task

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/yuedun/zhuque/db"
	"github.com/yuedun/zhuque/pkg/message"
	"github.com/yuedun/zhuque/pkg/user"
	"github.com/yuedun/zhuque/util"

	"github.com/gin-gonic/gin"
)

//List
func List(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset := (page - 1) * limit
	taskName := c.Query("searchParams[taskName]")
	projectName := c.Query("searchParams[projectName]")
	username := c.Query("searchParams[username]")
	var task Task
	task.TaskName = taskName
	task.Project = projectName
	task.Username = username
	serverService := NewService(db.DB)
	list, count, err := serverService.GetTaskList(offset, limit, task)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  list,
		"msg":   "ok",
	})
}

//WaitList 等待发布任务
func WaitList(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	from := c.Query("from")
	serverService := NewService(db.DB)
	list, err := serverService.WaitTaskList(from)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": len(list),
		"data":  list,
		"msg":   "ok",
	})
}

//GetTaskInfo
func GetTaskInfo(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.Param("id"))
	name := c.Param("name")
	taskService := NewService(db.DB)
	taskObj := Task{
		ID:       taskID,
		TaskName: name,
	}
	task, err := taskService.GetTaskInfo(taskObj)
	if err != nil {
		log.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    task,
		"message": "ok",
	})
}

//GetTaskInfoBySql
func GetTaskInfoBySql(c *gin.Context) {
	taskService := NewService(db.DB)
	task, err := taskService.GetTaskInfoBySQL()
	if err != nil {
		log.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": task,
	})
}

//CreateTask
func CreateTask(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	taskService := NewService(db.DB)
	task := Task{}
	if err := c.ShouldBind(&task); err != nil {
		panic(err)
	}
	task.CreatedAt = time.Now()
	_, err := taskService.CreateTask(&task)
	if err != nil {
		log.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    task,
		"message": "ok",
	})
}

//UpdateTask post json
func UpdateTask(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	taskID, _ := strconv.Atoi(c.Param("id"))
	taskService := NewService(db.DB)
	var task Task
	if err := c.ShouldBind(&task); err != nil {
		panic(err)
	}
	task.ID = taskID
	err := taskService.UpdateTask(taskID, &task)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    task,
		"message": "ok",
	})
}

//DeleteTask
func DeleteTask(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.Param("id"))
	taskService := NewService(db.DB)
	err := taskService.DeleteTask(taskID)
	if err != nil {
		log.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

// Approve 发布审批 发邮件，钉钉等消息
func Approve(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	taskService := NewService(db.DB)
	var params map[string]interface{}
	if err := c.ShouldBind(&params); err != nil {
		panic(err)
	}
	task, err := taskService.Approve(params)
	if err != nil {
		panic(err)
	}
	content := ""
	if params["nowRelease"] == true {
		content = fmt.Sprintf("【朱雀】发布单【%s】可立即发布", task.TaskName)
	}
	releaseSta, ok := params["releaseState"]
	if ok {
		rstate := releaseSta.(float64)
		if int(rstate) == 0 {
			content = fmt.Sprintf("【朱雀】发布单【%s】当前不可发布，原因：%s", task.TaskName, task.ApproveMsg)
		}
	}
	bodyObj := make(map[string]interface{})
	bodyObj["msgtype"] = "text"
	bodyObj["text"] = map[string]interface{}{
		"content": content,
	}
	// 发送给有项目权限的人
	userService := user.NewService(db.DB)
	mailTo, err := userService.GetProjectUsersEmail(task.Project)
	if err != nil {
		//邮件错误忽略，不影响主流程
		log.Println(err)
	}
	// mailTo := strings.Split(users, ";")
	messageService := message.NewMessage()
	// 异步发送，避免阻塞，发送成功与否都没关系
	go messageService.SendDingTalk(util.Conf.DingTalk, bodyObj)
	go messageService.SendEmail(task.TaskName, content, mailTo)
	c.JSON(http.StatusOK, gin.H{
		"data":    "",
		"message": "ok",
	})
}
