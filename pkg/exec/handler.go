package exec

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuedun/zhuque/db"
	"github.com/yuedun/zhuque/pkg/task"

	"os/exec"
)

// Send 发送命令到服务器
func Send(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	userCmd, ok := c.GetPostForm("cmd")
	if !ok || userCmd == "" {
		panic(errors.New("命令无效！"))
	}
	log.Println("用户输入命令：", userCmd)
	var cmdOut []byte
	var err error
	var cmd *exec.Cmd
	// 执行单个shell命令时, 直接运行即可
	cmd = exec.Command("bash", "-c", userCmd)
	if cmdOut, err = cmd.CombinedOutput(); err != nil {
		log.Println("输出错误：", err)
		log.Println("输出错误2：", string(cmdOut))
		c.JSON(200, gin.H{
			"message": err,
			"data":    strings.ReplaceAll(string(cmdOut), "\n", "<br>"),
		})
		return
	}
	// 默认输出有一个换行
	log.Println(string(cmdOut))

	c.JSON(200, gin.H{
		"message": "ok",
		"data":    strings.ReplaceAll(string(cmdOut), "\n", "<br>"),
	})
}

// Server 快捷发布，只需要选择项目即可发布
func Server(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	cmdParam, ok := c.GetPostForm("cmd")
	if !ok || cmdParam == "" {
		panic(errors.New("命令无效！"))
	}
	userID, ok := c.GetPostForm("userID")
	if !ok || userID == "" {
		panic(errors.New("用户ID无效！"))
	}
	username, ok := c.GetPostForm("username")
	if !ok || username == "" {
		panic(errors.New("用户名无效！"))
	}
	userCmd := fmt.Sprintf("pm2 deploy projects/%s/ecosystem.config.js production --force", cmdParam)
	log.Println("用户输入命令：", userCmd)
	var cmdOut []byte
	var err error
	var cmd *exec.Cmd
	// 执行单个shell命令时, 直接运行即可
	cmd = exec.Command("bash", "-c", userCmd)
	taskServer := task.NewService(db.SQLLite)
	task := task.Task{
		TaskName: "",
		Project:  cmdParam,
		UserID:   userID,
		Username: username,
	}

	if cmdOut, err = cmd.CombinedOutput(); err != nil {
		log.Println("输出错误：", err)
		log.Println("输出错误2：", string(cmdOut))
		task.ReleaseResult = 0       //失败
		taskServer.CreateTask(&task) //保存失败发布记录
		c.JSON(200, gin.H{
			"message": err,
			"data":    strings.ReplaceAll(string(cmdOut), "\n", "<br>"),
		})
		return
	}
	// 默认输出有一个换行
	log.Println(string(cmdOut))
	task.ReleaseResult = 1       //失败
	taskServer.CreateTask(&task) //保存成功发布记录
	c.JSON(200, gin.H{
		"message": "ok",
		"data":    strings.ReplaceAll(string(cmdOut), "\n", "<br>"),
	})
}
