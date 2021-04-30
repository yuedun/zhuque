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
	"github.com/yuedun/zhuque/pkg/project"
	"github.com/yuedun/zhuque/pkg/task"
	"github.com/yuedun/zhuque/pkg/user"
	"github.com/yuedun/zhuque/util"
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
	execService := NewService(db.DB)
	cmdOut, err := execService.CmdSync(userCmd)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err,
			"data":    strings.ReplaceAll(string(cmdOut), "\n", "<br>"),
		})
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
	projectName, ok := c.GetPostForm("cmd")
	if !ok || projectName == "" {
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
	projectServer := project.NewService(db.DB)
	projectObj := project.Project{
		Name: projectName,
	}
	project, err := projectServer.GetProjectInfo(projectObj)
	if err != nil {
		panic(err)
	}
	taskServer := task.NewService(db.DB)

	execService := NewService(db.DB)
	resCode := 1 // code=1是直接发布，code=2是审核发布
	resData := ""
	// 判断发布类型，pm2还是scp
	// scp发布类型
	if project.DeployType == "scp" {
		// 1.创建发布单
		task := task.Task{
			TaskName:     taskName,
			Project:      projectName,
			UserID:       userID,
			ReleaseState: 2, //待发布
			Username:     username,
			From:         "single",
		}
		taskID, err := taskServer.CreateTask(&task)
		if err != nil {
			panic(err)
		}
		// 如果是测服直接发布
		if util.Conf.Env == "prod" {
			execService.SendMessage(task)
			resCode = 2
			resData = fmt.Sprintf("%d分钟后可发布", util.Conf.DelayDeploy)
		} else {
			output, err := execService.DeployControl(project.ID, taskID)
			if err != nil {
				task.ReleaseState = 0
				taskServer.UpdateTask(taskID, &task)
				output = err.Error()
			} else {
				task.ReleaseState = 1
				taskServer.UpdateTask(taskID, &task)
			}
			resData = strings.ReplaceAll(string(output), "\n", "<br>")
		}
	} else {
		userCmd := fmt.Sprintf("pm2 deploy projects/%s/ecosystem.config.js production --force", projectName)
		if restart == "on" {
			// 由于pm2的项目名和管理的项目名不能完全保持一致，所以如果一个pm2下跑多个服务都只能重启，但是reload可以实现不停服重启
			userCmd = fmt.Sprintf("pm2 deploy projects/%s/ecosystem.config.js production exec 'git pull && pm2 reload ecosystem.config.js' --force && pm2 list", projectName)
		}
		log.Println("用户输入命令：", userCmd)

		// 1.创建发布单
		task := task.Task{
			TaskName:     taskName,
			Project:      projectName,
			UserID:       userID,
			ReleaseState: 2, //待发布
			Username:     username,
			Cmd:          userCmd,
			From:         "single",
		}
		taskID, err := taskServer.CreateTask(&task)
		if err != nil {
			panic(err)
		}
		if util.Conf.Env == "prod" {
			// 生产n分钟后发布
			execService.SendMessage(task)
			resCode = 2
			resData = fmt.Sprintf("%d分钟后可发布", util.Conf.DelayDeploy)
		} else {
			// 测服直接发布 resCode=1，前端调用发布接口

			// resData, err = taskServer.ReleaseTask(taskID)
			resData = fmt.Sprint(taskID)
		}
	}
	c.JSON(200, gin.H{
		"code":    resCode, //code=1是直接发布，code=2是审核发布
		"message": err,
		"data":    resData,
	})
}

// 创建发布单-pm2发布模式
func CreateTaskForPM2(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	//cmd的值是项目名
	projectName, ok := c.GetPostForm("cmd")
	if !ok || projectName == "" {
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

	taskServer := task.NewService(db.DB)
	execService := NewService(db.DB)

	resCode := 1 // code=1是直接发布，code=2是审核发布
	resData := ""

	userCmd := fmt.Sprintf("pm2 deploy projects/%s/ecosystem.config.js production --force", projectName)
	if restart == "on" {
		// 由于pm2的项目名和管理的项目名不能完全保持一致，所以如果一个pm2下跑多个服务都只能重启，但是reload可以实现不停服重启
		userCmd = fmt.Sprintf("pm2 deploy projects/%s/ecosystem.config.js production exec 'git pull && pm2 reload ecosystem.config.js' --force && pm2 list", projectName)
	}
	log.Println("用户输入命令：", userCmd)

	// 1.创建发布单
	task := task.Task{
		TaskName:     taskName,
		Project:      projectName,
		UserID:       userID,
		ReleaseState: 2, //待发布
		Username:     username,
		Cmd:          userCmd,
		From:         "single",
	}
	taskID, err := taskServer.CreateTask(&task)
	if err != nil {
		panic(err)
	}
	if util.Conf.Env == "prod" {
		// 生产n分钟后发布
		execService.SendMessage(task)
		resCode = 2
		resData = fmt.Sprintf("%d分钟后可发布", util.Conf.DelayDeploy)
	} else {
		// 测服直接发布 resCode=1，前端调用发布接口

		// resData, err = taskServer.ReleaseTask(taskID)
		resData = fmt.Sprint(taskID)
	}

	c.JSON(200, gin.H{
		"code":    resCode, //code=1是直接发布，code=2是审核发布
		"message": err,
		"data":    resData,
	})
}

// 创建发布单-scp发布模式
func CreateTaskForSCP(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	//cmd的值是项目名
	projectName, ok := c.GetPostForm("cmd")
	if !ok || projectName == "" {
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

	taskServer := task.NewService(db.DB)
	execService := NewService(db.DB)

	resCode := 1 // code=1是直接发布，code=2是审核发布
	resData := ""

	userCmd := fmt.Sprintf("pm2 deploy projects/%s/ecosystem.config.js production --force", projectName)
	if restart == "on" {
		// 由于pm2的项目名和管理的项目名不能完全保持一致，所以如果一个pm2下跑多个服务都只能重启，但是reload可以实现不停服重启
		userCmd = fmt.Sprintf("pm2 deploy projects/%s/ecosystem.config.js production exec 'git pull && pm2 reload ecosystem.config.js' --force && pm2 list", projectName)
	}
	log.Println("用户输入命令：", userCmd)

	// 1.创建发布单
	task := task.Task{
		TaskName:     taskName,
		Project:      projectName,
		UserID:       userID,
		ReleaseState: 2, //待发布
		Username:     username,
		Cmd:          userCmd,
		From:         "single",
	}
	taskID, err := taskServer.CreateTask(&task)
	if err != nil {
		panic(err)
	}
	if util.Conf.Env == "prod" {
		// 生产n分钟后发布
		execService.SendMessage(task)
		resCode = 2
		resData = fmt.Sprintf("%d分钟后可发布", util.Conf.DelayDeploy)
	} else {
		// 测服直接发布 resCode=1，前端调用发布接口

		// resData, err = taskServer.ReleaseTask(taskID)
		resData = fmt.Sprint(taskID)
	}

	c.JSON(200, gin.H{
		"code":    resCode, //code=1是直接发布，code=2是审核发布
		"message": err,
		"data":    resData,
	})
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
		userCmd = "pm2 deploy projects/%s/ecosystem.config.js production exec 'git pull && pm2 reload ecosystem.config.js' --force && pm2 list"
	}
	log.Println("用户输入命令：", userCmd)

	// 1.创建发布单
	taskServer := task.NewService(db.DB)
	task := task.Task{
		TaskName:     taskName,
		Project:      strings.Join(projects, ","),
		UserID:       userID,
		Username:     username,
		ReleaseState: 2,
		Cmd:          userCmd,
		From:         "multi",
	}
	_, err := taskServer.CreateTask(&task)
	if err != nil {
		panic(err)
	}

	// 如果是测服直接发布
	var cmdOut string
	if util.Conf.Env == "prod" {
		// 发送消息通知
		content := fmt.Sprintf("【朱雀】发布单【%s】将在%d分钟后发布%s。提交人：%s", task.TaskName, util.Conf.DelayDeploy, task.Project, task.Username)
		log.Printf(content)
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
		c.JSON(200, gin.H{
			"code":    2, //code=1是直接发布，code=2是审核发布
			"message": "ok",
			"data":    fmt.Sprintf("%d分钟后可发布", util.Conf.DelayDeploy),
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

// 创建发布单-pm2发布模式，发布多个项目，主要是正产环境集群发布
func CreateTaskForPM2V2(c *gin.Context) {
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

	resCode := 1 // code=1是直接发布，code=2是审核发布
	resData := ""

	userCmd := "pm2 deploy projects/%s/ecosystem.config.js production --force"
	if restart == "on" {
		// 由于pm2的项目名和管理的项目名不能完全保持一致，所以如果一个pm2下跑多个服务都只能重启，但是reload可以实现不停服重启
		userCmd = "pm2 deploy projects/%s/ecosystem.config.js production exec 'git pull && pm2 reload ecosystem.config.js' --force && pm2 list"
	}
	log.Println("用户输入命令：", userCmd)

	// 1.创建发布单
	taskServer := task.NewService(db.DB)
	task := task.Task{
		TaskName:     taskName,
		Project:      strings.Join(projects, ","),
		UserID:       userID,
		Username:     username,
		ReleaseState: 2,
		Cmd:          userCmd,
		From:         "multi",
	}
	taskID, err := taskServer.CreateTask(&task)
	if err != nil {
		panic(err)
	}

	// 如果是测服直接发布
	if util.Conf.Env == "prod" {
		// 发送消息通知
		content := fmt.Sprintf("【朱雀】发布单【%s】将在%d分钟后发布%s。提交人：%s", task.TaskName, util.Conf.DelayDeploy, task.Project, task.Username)
		log.Printf(content)
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

		resCode = 2
		resData = fmt.Sprintf("%d分钟后可发布", util.Conf.DelayDeploy)
	} else {
		resCode = 1
		resData = fmt.Sprint(taskID)
	}

	c.JSON(200, gin.H{
		"code":    resCode, //code=1是直接发布，code=2是审核发布
		"message": err,
		"data":    resData,
	})
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
	taskServer := task.NewService(db.DB)
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
	taskServer := task.NewService(db.DB)
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
