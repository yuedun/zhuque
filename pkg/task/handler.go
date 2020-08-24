package task

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/yuedun/zhuque/db"

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

//GetTaskInfo
func GetTaskInfo(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.Param("id"))
	name := c.Param("name")
	ip := c.Param("ip")
	taskService := NewService(db.SQLLite)
	taskObj := Task{
		ID:       taskID,
		TaskName: name,
		Status:   ip,
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
	taskService := NewService(db.SQLLite)
	var task Task
	taskID, _ := strconv.Atoi(c.Param("id"))
	//task.Addr = c.PostForm("addr")
	if err := c.ShouldBind(&task); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data":    nil,
			"message": "err",
		})
	} else {
		err := taskService.UpdateTask(taskID, &task)
		if err != nil {
			fmt.Println("err:", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"data":    task,
			"message": "ok",
		})
	}
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
