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
	userID, _ := strconv.Atoi(c.Param("id"))
	name := c.Param("name")
	ip := c.Param("ip")
	userService := NewService(db.SQLLite)
	userObj := Task{
		ID:       userID,
		TaskName: name,
		Status:   ip,
	}
	user, err := userService.GetTaskInfo(userObj)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "ok",
	})
}

//GetTaskInfoBySql
func GetTaskInfoBySql(c *gin.Context) {
	userService := NewService(db.SQLLite)
	user, err := userService.GetTaskInfoBySQL()
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": user,
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
	userService := NewService(db.SQLLite)
	user := Task{}
	if err := c.ShouldBind(&user); err != nil {
		panic(err)
	}
	user.CreatedAt = time.Now()
	err := userService.CreateTask(&user)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "ok",
	})
}

//UpdateTask post json
func UpdateTask(c *gin.Context) {
	userService := NewService(db.SQLLite)
	var user Task
	userID, _ := strconv.Atoi(c.Param("id"))
	//user.Addr = c.PostForm("addr")
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data":    nil,
			"message": "err",
		})
	} else {
		err := userService.UpdateTask(userID, &user)
		if err != nil {
			fmt.Println("err:", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"data":    user,
			"message": "ok",
		})
	}
}

//DeleteTask
func DeleteTask(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))
	userService := NewService(db.SQLLite)
	err := userService.DeleteTask(userID)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
