package task

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/yuedun/zhuque/db"
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
	var task Task
	serverService := NewService(db.SQLLite)
	list, err := serverService.GetTaskList(task)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": list,
		"msg":  "ok",
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
	task := Task{ReleaseState: 2}
	serverService := NewService(db.SQLLite)
	list, err := serverService.GetTaskList(task)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": list,
		"msg":  "ok",
	})
}

//GetTaskInfo
func GetTaskInfo(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.Param("id"))
	name := c.Param("name")
	taskService := NewService(db.SQLLite)
	taskObj := Task{
		ID:       taskID,
		TaskName: name,
	}
	task, err := taskService.GetTaskInfo(taskObj)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    task,
		"message": "ok",
	})
}

//GetTaskInfoBySql
func GetTaskInfoBySql(c *gin.Context) {
	taskService := NewService(db.SQLLite)
	task, err := taskService.GetTaskInfoBySQL()
	if err != nil {
		fmt.Println("err:", err)
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
	taskService := NewService(db.SQLLite)
	task := Task{}
	if err := c.ShouldBind(&task); err != nil {
		panic(err)
	}
	task.CreatedAt = time.Now()
	err := taskService.CreateTask(&task)
	if err != nil {
		fmt.Println("err:", err)
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
	taskService := NewService(db.SQLLite)
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
	taskService := NewService(db.SQLLite)
	err := taskService.DeleteTask(taskID)
	if err != nil {
		fmt.Println("err:", err)
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
	taskService := NewService(db.SQLLite)
	var params map[string]interface{}
	if err := c.ShouldBind(&params); err != nil {
		panic(err)
	}
	log.Println(">>>>>>>>>>>>", params)
	task, err := taskService.Approve(params)
	if err != nil {
		panic(err)
	}
	content := ""
	if params["nowRelease"] == true {
		content = fmt.Sprintf("【朱雀】发布单【%s】可立即发布", task.TaskName)
	}
	releaseSta := params["releaseState"].(float64)
	if int(releaseSta) == 0 {
		content = fmt.Sprintf("【朱雀】发布单【%s】当前不可发布，原因：%s", task.TaskName, task.ApproveMsg)
	}
	bodyObj := make(map[string]interface{})
	bodyObj["msgtype"] = "text"
	bodyObj["text"] = map[string]interface{}{
		"content": content,
	}
	util.SendDingTalk(util.Conf.DingTalk, bodyObj)
	c.JSON(http.StatusOK, gin.H{
		"data":    "",
		"message": "ok",
	})
}
