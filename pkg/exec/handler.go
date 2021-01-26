package exec

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
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
	cmdOut, err := ExecCmdSync(userCmd)
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
		userCmd = fmt.Sprintf("pm2 deploy projects/%s/ecosystem.config.js production exec 'git pull && pm2 reload ecosystem.config.js' --force", cmdParam)
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
		From:         "single",
	}
	err := taskServer.CreateTask(&task)
	if err != nil {
		panic(err)
	}

	// 如果是测服直接发布
	var cmdOut string
	if util.Conf.Env == "prod" {
		// content消息内容
		content := fmt.Sprintf("【朱雀】发布单【%s】将在%d分钟后发布%s。提交人：%s", task.TaskName, util.Conf.DelayDeploy, task.Project, task.Username)

		//bodyObj 钉钉消息体
		bodyObj := make(map[string]interface{})
		bodyObj["msgtype"] = "text"
		bodyObj["text"] = map[string]interface{}{
			"content": content,
		}
		// 发送给有项目权限的人
		userService := user.NewService(db.SQLLite)
		mailTo, err := userService.GetProjectUsersEmail(task.Project)
		if err != nil {
			//邮件错误忽略，不影响主流程
			log.Println(err)
		}
		// mailTo := strings.Split(users, ";")
		messageService := message.NewMessage()
		// 异步发送，避免阻塞，发送成功与否都没关系
		go messageService.SendDingTalk(util.Conf.DingTalk, bodyObj)
		go messageService.SendEmailV2(task.TaskName, content, mailTo)
		c.JSON(200, gin.H{
			"code":    2, // code=1是直接发布，code=2是审核发布
			"message": "ok",
			"data":    fmt.Sprintf("%d分钟后可发布", util.Conf.DelayDeploy),
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
		userCmd = "pm2 deploy projects/%s/ecosystem.config.js production exec 'git pull && pm2 reload ecosystem.config.js' --force"
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
		From:         "multi",
	}
	err := taskServer.CreateTask(&task)
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
		userService := user.NewService(db.SQLLite)
		mailTo, err := userService.GetProjectUsersEmail(task.Project)
		if err != nil {
			//邮件错误忽略，不影响主流程
			log.Println(err)
		}
		// mailTo := strings.Split(users, ";")
		messageService := message.NewMessage()
		// 异步发送，避免阻塞，发送成功与否都没关系
		go messageService.SendDingTalk(util.Conf.DingTalk, bodyObj)
		go messageService.SendEmailV2(task.TaskName, content, mailTo)
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

// Deploy 第一次部署
// 1.git拉代码到发布机指定目录下，项目名和发布系统中保持一致。
// 2.安装依赖，编译等操作
// 3.同步代码到目标机器的指定目录下
func DeployProject(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	projectID, _ := strconv.Atoi(c.Param("id"))
	projectObj := project.Project{
		ID: projectID,
	}
	projectService := project.NewService(db.SQLLite)
	projectResult, _ := projectService.GetProjectInfo(projectObj)
	var config map[string]interface{}
	err := json.Unmarshal([]byte(projectResult.Config), &config)
	if err != nil {
		panic(err)
	}
	production := config["deploy"].(map[string]interface{})
	productionJSON, err := json.Marshal(production["production"])
	log.Println(string(productionJSON))
	var deployConfig project.DeployConfig
	err = json.Unmarshal(productionJSON, &deployConfig)
	if err != nil {
		panic(err)
	}
	log.Println(deployConfig.Host)
	var buffer bytes.Buffer
	// 目录不存在则git clone
	if ex := util.PathExists(util.Conf.APPDir); ex == false {
		log.Println(">>>>>>>>>>>", ex)
		cmd1 := fmt.Sprintf("git clone -b %s %s %s", deployConfig.Ref, deployConfig.Repo, path.Join(util.Conf.APPDir, projectResult.Name))
		log.Println("第一步：检出代码：", cmd1)
		cmdOut, err := ExecCmdSync(cmd1)
		if err != nil {
			panic(err)
		}
		buffer.Write(cmdOut)
	}
	cmd2 := "cd " + path.Join(util.Conf.APPDir, projectResult.Name) + " && npm i"
	log.Println("第二步：安装依赖：", cmd2)
	cmdOut, err := ExecCmdSync(cmd2)
	if err != nil {
		panic(err)
	}
	buffer.Write(cmdOut)
	cmd3 := fmt.Sprintf("rsync -av %s %s", path.Join(util.Conf.APPDir, projectResult.Name), fmt.Sprintf("%s@%s:%s", deployConfig.User, deployConfig.Host, deployConfig.Path))
	log.Println("第三步：同步代码：", cmd3)
	cmdOut, err = ExecCmdSync(cmd3)
	if err != nil {
		panic(err)
	}
	// 默认输出有一个换行
	log.Println("执行结果：", string(cmdOut))
	c.JSON(200, gin.H{
		"data": string(buffer.Bytes()),
	})
}
