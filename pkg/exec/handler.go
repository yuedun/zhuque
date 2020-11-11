package exec

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuedun/zhuque/db"
	"github.com/yuedun/zhuque/pkg/message"
	"github.com/yuedun/zhuque/pkg/task"
	"github.com/yuedun/zhuque/util"

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
	//cmd的值是项目名
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
	taskName, ok := c.GetPostForm("taskName")
	if !ok || taskName == "" {
		panic(errors.New("用户名无效！"))
	}
	restart, ok := c.GetPostForm("restart")
	userCmd := fmt.Sprintf("pm2 deploy projects/%s/ecosystem.config.js production --force", cmdParam)
	if restart == "on" {
		// 由于pm2的项目名和管理的项目名不能完全保持一致，所以如果一个pm2下跑多个服务都只能重启，但是reload可以实现不停服重启
		userCmd = fmt.Sprintf("pm2 deploy projects/%s/ecosystem.config.js production exec 'git pull && pm2 reload all' --force", cmdParam)
	}
	log.Println("用户输入命令：", userCmd)

	// 1.创建发布单
	taskServer := task.NewService(db.SQLLite)
	task := task.Task{
		TaskName:     taskName,
		Project:      cmdParam,
		UserID:       userID,
		Username:     username,
		ReleaseState: 2,
		Cmd:          userCmd,
	}
	err := taskServer.CreateTask(&task)
	if err != nil {
		panic(err)
	}

	// 如果是测服直接发布
	var cmdOut string
	if util.Conf.Env == "prod" {
		// 发送消息通知
		content := fmt.Sprintf("【朱雀】发布单【%s】将在10分钟后发布。提交人：%s", task.TaskName, task.Username)
		log.Printf(content)
		bodyObj := make(map[string]interface{})
		bodyObj["msgtype"] = "text"
		bodyObj["text"] = map[string]interface{}{
			"content": content,
		}
		messageService := message.NewMessage()
		// 异步发送，避免阻塞，发送成功与否都没关系
		go messageService.SendDingTalk(util.Conf.DingTalk, bodyObj)
		go messageService.SendEmail(task.TaskName, content, util.Conf.EmailTo)
		c.JSON(200, gin.H{
			"code":    2, //code=1是直接发布，code=2是审核发布
			"message": "ok",
			"data":    "5分钟后可发布",
		})
	} else {
		cmdOut, err = taskServer.ReleaseTask(task.ID)
		c.JSON(200, gin.H{
			"code":    1, //code=1是直接发布，code=2是审核发布
			"message": err,
			"data":    cmdOut,
		})
	}
}

// ServerV2 快捷发布，发布多个项目，主要是正产环境集群发布
func ServerV2(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	//projects的值是项目名
	projects, ok := c.GetPostFormArray("projects[]")
	log.Println(">>>>>>>发布项目", projects)
	if !ok {
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
	taskName, ok := c.GetPostForm("taskName")
	if !ok || taskName == "" {
		panic(errors.New("用户名无效！"))
	}
	restart, ok := c.GetPostForm("restart")
	userCmd := "pm2 deploy projects/%s/ecosystem.config.js production --force"
	if restart == "on" {
		// 由于pm2的项目名和管理的项目名不能完全保持一致，所以如果一个pm2下跑多个服务都只能重启，但是reload可以实现不停服重启
		userCmd = "pm2 deploy projects/%s/ecosystem.config.js production exec 'git pull && pm2 reload all' --force"
	}
	log.Println("用户输入命令：", userCmd)

	// 1.创建发布单
	taskServer := task.NewService(db.SQLLite)
	task := task.Task{
		TaskName:     taskName,
		Project:      strings.Join(projects, ","),
		UserID:       userID,
		Username:     username,
		ReleaseState: 2,
		Cmd:          userCmd,
	}
	err := taskServer.CreateTask(&task)
	if err != nil {
		panic(err)
	}

	// 如果是测服直接发布
	var cmdOut string
	if util.Conf.Env == "prod" {
		// 发送消息通知
		content := fmt.Sprintf("【朱雀】发布单【%s】将在10分钟后发布", task.TaskName)
		log.Printf(content)
		bodyObj := make(map[string]interface{})
		bodyObj["msgtype"] = "text"
		bodyObj["text"] = map[string]interface{}{
			"content": content,
		}
		messageService := message.NewMessage()
		// 异步发送，避免阻塞，发送成功与否都没关系
		go messageService.SendDingTalk(util.Conf.DingTalk, bodyObj)
		go messageService.SendEmail(task.TaskName, content, util.Conf.EmailTo)
		c.JSON(200, gin.H{
			"code":    2, //code=1是直接发布，code=2是审核发布
			"message": "ok",
			"data":    "10分钟后可发布",
		})
	} else {
		cmdOut, err = taskServer.ReleaseTaskV2(task.ID)
		c.JSON(200, gin.H{
			"code":    1, //code=1是直接发布，code=2是审核发布
			"message": err,
			"data":    cmdOut,
		})
	}
}

// Release 发布操作，提交后的延时发布
func Release(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	taskID, _ := strconv.Atoi(c.Param("id"))
	taskServer := task.NewService(db.SQLLite)
	// TODO需要验证是否可发布
	cmdOut, err := taskServer.ReleaseTask(taskID)
	if err != nil {
		panic(err)
	}
	c.JSON(200, gin.H{
		"message": err,
		"data":    cmdOut,
	})
}

// ReleaseV2 发布操作，提交后的延时发布
func ReleaseV2(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	taskID, _ := strconv.Atoi(c.Param("id"))
	taskServer := task.NewService(db.SQLLite)
	// TODO需要验证是否可发布
	cmdOut, err := taskServer.ReleaseTaskV2(taskID)
	if err != nil {
		panic(err)
	}
	c.JSON(200, gin.H{
		"message": err,
		"data":    cmdOut,
	})
}
