package exec

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
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
	projectName, ok := c.GetPostForm("project")
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

	taskServer := task.NewService(db.DB)
	execService := NewService(db.DB)

	resCode := 1 // code=1是直接发布，code=2是审核发布
	resData := ""

	// restart, ok := c.GetPostForm("restart")
	// userCmd := fmt.Sprintf("pm2 deploy projects/%s/ecosystem.config.js production --force", projectName)
	// if restart == "on" {
	// 	// 由于pm2的项目名和管理的项目名不能完全保持一致，所以如果一个pm2下跑多个服务都只能重启，但是reload可以实现不停服重启
	// 	userCmd = fmt.Sprintf("pm2 deploy projects/%s/ecosystem.config.js production exec 'git pull && pm2 reload ecosystem.config.js' --force && pm2 list", projectName)
	// }

	// 1.创建发布单
	task := task.Task{
		TaskName:     taskName,
		Project:      projectName,
		UserID:       userID,
		ReleaseState: task.Ready, //待发布
		Username:     username,
		// Cmd:          userCmd,
		From:       "single",
		DeployType: "pm2",
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
	cmdOut, err := taskServer.ReleaseTask(taskID)
	if err != nil {
		panic(err)
	}
	c.JSON(200, gin.H{
		"message": err,
		"data":    cmdOut,
	})
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
	// restart, ok := c.GetPostForm("restart")
	// userCmd := "pm2 deploy projects/%s/ecosystem.config.js production --force"
	// if restart == "on" {
	// 	// 由于pm2的项目名和管理的项目名不能完全保持一致，所以如果一个pm2下跑多个服务都只能重启，但是reload可以实现不停服重启
	// 	userCmd = "pm2 deploy projects/%s/ecosystem.config.js production exec 'git pull && pm2 reload ecosystem.config.js' --force && pm2 list"
	// }

	resCode := 1 // code=1是直接发布，code=2是审核发布
	resData := ""

	// 1.创建发布单
	taskServer := task.NewService(db.DB)
	task := task.Task{
		TaskName:     taskName,
		Project:      strings.Join(projects, ","),
		UserID:       userID,
		Username:     username,
		ReleaseState: task.Ready,
		// Cmd:          userCmd,
		From:       "multi",
		DeployType: "pm2",
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
	cmdOut, err := taskServer.ReleaseTaskV2(taskID)
	if err != nil {
		panic(err)
	}
	c.JSON(200, gin.H{
		"message": err,
		"data":    cmdOut,
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
	projectName, ok := c.GetPostForm("project")
	if !ok || projectName == "" {
		panic(errors.New("项目名无效！"))
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
	// restart, ok := c.GetPostForm("restart")

	taskServer := task.NewService(db.DB)
	execService := NewService(db.DB)

	resCode := 1 // code=1是直接发布，code=2是审核发布
	resData := ""

	// scp发布类型
	// 1.创建发布单
	task := task.Task{
		TaskName:     taskName,
		Project:      projectName,
		UserID:       userID,
		ReleaseState: task.Ready, //待发布
		Username:     username,
		From:         "single",
		DeployType:   "scp",
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
		// 测服直接发布 resCode=1，前端调用发布接口
		resCode = 1
		resData = fmt.Sprint(taskID)
	}

	c.JSON(200, gin.H{
		"code":    resCode, //code=1是直接发布，code=2是审核发布
		"message": err,
		"data":    resData,
	})
}

// ReleaseForScp 发布操作，提交后的延时发布
func ReleaseForSCP(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	resCode := 0
	resData := ""
	taskID, _ := strconv.Atoi(c.Param("id"))
	taskServer := task.NewService(db.DB)
	execService := NewService(db.DB)
	projectServer := project.NewService(db.DB)

	taskObj := task.Task{
		ID: taskID,
	}
	taskRecord, err := taskServer.GetTaskInfo(taskObj)
	if err != nil {
		panic(err)
	}
	taskRecord.ReleaseState = task.Releaseing
	taskServer.UpdateTask(taskID, &taskRecord)

	projectObj := project.Project{
		Name: taskRecord.Project,
	}
	project, err := projectServer.GetProjectInfo(projectObj)
	if err != nil {
		panic(err)
	}
	claims := jwt.ExtractClaims(c)
	userID64 := claims["user_id"].(float64)
	userID := int(userID64)
	output, err := execService.DeployControl(project, taskID, strconv.Itoa(userID))
	if err != nil {
		resCode = 1
		taskRecord.ReleaseState = task.Fail //失败
	} else {
		taskRecord.ReleaseState = task.Success //成功
	}
	taskServer.UpdateTask(taskID, &taskRecord)
	resData = strings.ReplaceAll(string(output), "\n", "<br>")
	c.JSON(200, gin.H{
		"code":    resCode,
		"message": err,
		"data":    resData,
	})
}
