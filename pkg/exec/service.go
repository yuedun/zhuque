package exec

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path"

	"github.com/jinzhu/gorm"
	"github.com/yuedun/zhuque/db"
	"github.com/yuedun/zhuque/pkg/message"
	"github.com/yuedun/zhuque/pkg/project"
	"github.com/yuedun/zhuque/pkg/task"
	"github.com/yuedun/zhuque/pkg/user"
	"github.com/yuedun/zhuque/util"
)

type (
	ExecService interface {
		DeployControl(project project.Project, taskID int) ([]byte, error)
		CmdSync(userCmd string) ([]byte, error)
		CloneRepo(deployConfig project.DeployConfig, projectName string) ([]byte, error)
		GitPull(deployConfig project.DeployConfig, projectName string) ([]byte, error)
		PreBuild(deployConfig project.DeployConfig, projectName string) ([]byte, error)
		Build(deployConfig project.DeployConfig, projectName string) ([]byte, error)
		SyncCode(deployConfig project.DeployConfig, projectName string) ([]byte, error)
		PreDeploy(deployConfig project.DeployConfig, projectName string) ([]byte, error)
		PostDeploy(deployConfig project.DeployConfig, projectName string) ([]byte, error)
		SendMessage(task task.Task)
	}
)
type execService struct {
	db *gorm.DB
}

/*NewService 初始化结构体*/
func NewService(db *gorm.DB) ExecService {
	return &execService{
		db: db,
	}
}

// DeployControl 发布流程控制
func (u *execService) DeployControl(projectObj project.Project, taskID int) ([]byte, error) {
	var deployConfig project.DeployConfig
	var buffer bytes.Buffer
	var output []byte
	err := json.Unmarshal([]byte(projectObj.Config), &deployConfig)
	if err != nil {
		log.Println("项目配置解析失败，请检查配置json是否正确1:", err)
		buffer.Write([]byte(err.Error()))
		return buffer.Bytes(), err
	}
	if deployConfig.User == "" || len(deployConfig.Host) == 0 || deployConfig.Ref == "" || deployConfig.Repo == "" || deployConfig.Path == "" {
		log.Println("请检查配置是否完整")
		buffer.Write([]byte("请检查配置是否完整"))
		return buffer.Bytes(), errors.New("请检查配置是否完整")
	}

	// 1.克隆代码
	projectDirName := fmt.Sprintf("%s", projectObj.Name)
	log.Println("克隆项目名：", projectDirName)
	if exists := util.PathExists(path.Join(util.Conf.APPDir, projectDirName)); exists == false {
		// 分支，gitrepo，
		output, err = u.CloneRepo(deployConfig, projectDirName)
		if err != nil {
			buffer.Write([]byte(err.Error()))
			return buffer.Bytes(), err
		}
		buffer.Write(output)
	} else {
		buffer.Write([]byte("项目已存在，跳过克隆代码。\n"))
	}

	// 1.1拉新代码
	output, err = u.GitPull(deployConfig, projectDirName)
	if err != nil {
		buffer.Write([]byte(err.Error()))
		return buffer.Bytes(), err
	}
	buffer.Write(output)

	// 2.编译前置
	if deployConfig.PreBuild != "" {
		output, err = u.PreBuild(deployConfig, projectDirName)
		if err != nil {
			buffer.Write([]byte(err.Error()))
			return buffer.Bytes(), err
		}
		buffer.Write(output)
	}
	// 3.编译
	if deployConfig.Build != "" {
		output, err = u.Build(deployConfig, projectDirName)
		if err != nil {
			buffer.Write([]byte(err.Error()))
			return buffer.Bytes(), err
		}
		buffer.Write(output)
	}

	// 4.同步代码到远程服务器 发生错误停止往下执行
	output, err = u.SyncCode(deployConfig, projectDirName)
	if err != nil {
		buffer.Write([]byte(err.Error()))
		return buffer.Bytes(), err
	}
	buffer.Write(output)

	// 5.重启服务前置
	if deployConfig.PreDeploy != "" {
		output, err = u.PreDeploy(deployConfig, projectDirName)
		if err != nil {
			buffer.Write([]byte(err.Error()))
			return buffer.Bytes(), err
		}
		buffer.Write(output)
	}

	// 6.同步代码到远程应用服务器后执行命令，如重启
	if deployConfig.PostDeploy != "" {
		output, err = u.PostDeploy(deployConfig, projectDirName)
		if err != nil {
			log.Println("远程命令执行异常：", err)
			buffer.Write([]byte(err.Error()))
			return buffer.Bytes(), err
		}
		buffer.Write(output)
	}
	return buffer.Bytes(), nil
}

// CmdSync 同步执行命令
func (u *execService) CmdSync(userCmd string) ([]byte, error) {
	var stdoutStderr []byte
	var cmd *exec.Cmd
	// 执行单个shell命令时, 直接运行即可
	log.Println("执行命令：", userCmd)
	cmd = exec.Command("bash", "-c", userCmd)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("错误码：", err)
		return nil, errors.New(string(stdoutStderr))
	}
	log.Println(userCmd, "执行结果：\n", string(stdoutStderr))
	return stdoutStderr, nil
}

// CloneRepo clone代码
func (u *execService) CloneRepo(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
	//分支，gitrepo，目录名
	cmd1 := fmt.Sprintf("git clone -b %s %s %s", deployConfig.Ref, deployConfig.Repo, path.Join(util.Conf.APPDir, projectName))
	log.Println("第一步：检出代码：", cmd1)
	cmdOut, err := u.CmdSync(cmd1)
	if err != nil {
		log.Println("第一步：检出代码执行失败：", err)
		return nil, err
	}
	return cmdOut, nil
}

// GitPull
func (u *execService) GitPull(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
	gitpull := fmt.Sprintf("cd %s; git pull origin %s; git log --oneline -1", path.Join(util.Conf.APPDir, projectName), deployConfig.Ref)
	cmdOut, err := u.CmdSync(gitpull)
	if err != nil {
		log.Println("拉取代码失败：", err)
		return nil, err
	}
	return cmdOut, nil
}

// PreBuild
func (u *execService) PreBuild(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
	build := fmt.Sprintf("cd %s; %s", path.Join(util.Conf.APPDir, projectName), deployConfig.PreBuild)
	cmdOut, err := u.CmdSync(build)
	if err != nil {
		log.Println("编译前置操作失败：", err)
		return nil, err
	}
	return cmdOut, nil
}

// Build
func (u *execService) Build(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
	build := fmt.Sprintf("cd %s; %s", path.Join(util.Conf.APPDir, projectName), deployConfig.Build)
	cmdOut, err := u.CmdSync(build)
	if err != nil {
		log.Println("编译操作失败：", err)
		return nil, err
	}
	return cmdOut, nil
}

// SyncCode 同步代码
func (u *execService) SyncCode(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
	hostLen := len(deployConfig.Host)
	ch := make(chan []byte, hostLen)
	errCh := make(chan error)
	for _, host := range deployConfig.Host {
		go func(host string, ch chan []byte, errCh chan error) {
			// 用户名@IP:远程目录
			remotePath := fmt.Sprintf("%s@%s:%s", deployConfig.User, host, deployConfig.Path)
			// rsync参数，宿主机项目，目标机地址
			cmd3 := fmt.Sprintf("rsync -av %s %s %s", deployConfig.RsyncArgs, path.Join(util.Conf.APPDir, projectName), remotePath)
			log.Println("同步代码：", cmd3)
			cmdput, err := u.CmdSync(cmd3)
			if err != nil {
				log.Println("同步代码执行失败：", err)
				errCh <- err
				return
			}
			ch <- cmdput
		}(host, ch, errCh)
	}
	i := 0
	var buffer bytes.Buffer
	var outErr error
	for {
		select {
		case out := <-ch:
			i++
			buffer.Write(out)

		case outErr = <-errCh:
			goto L
		}
		if i == hostLen {
			goto L
		}
	}
L:
	log.Println("SyncCode out")

	return buffer.Bytes(), outErr
}

// PreDeploy
func (u *execService) PreDeploy(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
	build := fmt.Sprintf("cd %s; %s", path.Join(util.Conf.APPDir, projectName), deployConfig.PreDeploy)
	cmdOut, err := u.CmdSync(build)
	if err != nil {
		log.Println("应用服务器重启前置：", err)
		return nil, err
	}
	return cmdOut, nil
}

// PostDeploy 远程应用服务器上发布 用户自定义命令：比如重启服务
// 利用ssh在远程服务器上执行
func (u *execService) PostDeploy(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
	hostLen := len(deployConfig.Host)
	ch := make(chan []byte, hostLen)
	// errch := make(chan error) //不需要缓冲，只要接收到一个错误就退出
	for _, host := range deployConfig.Host {
		go func(host string, ch chan []byte) {
			// 用户名@IP 命令
			ssh := fmt.Sprintf("ssh %s@%s \"cd %s; %s\"", deployConfig.User, host, path.Join(deployConfig.Path, projectName), deployConfig.PostDeploy)
			cmdput, err := u.CmdSync(ssh)
			if err != nil {
				log.Println("postDeploy过程 远程命令执行失败：", err)
				ch <- []byte(err.Error())
				return
			}
			ch <- cmdput
		}(host, ch)
	}
	i := 0
	var buffer bytes.Buffer
	for {
		select {
		case out := <-ch:
			i++
			buffer.Write(out)
		}
		if i == hostLen {
			goto L
		}
	}
L:
	log.Println("PostDeploy out")

	return buffer.Bytes(), nil
}

// SendMessage 生产发布消息通知
func (u *execService) SendMessage(task task.Task) {
	// content消息内容
	content := fmt.Sprintf("【朱雀】发布单【%s】将在%d分钟后发布%s。提交人：%s", task.TaskName, util.Conf.DelayDeploy, task.Project, task.Username)
	//bodyObj 钉钉消息体
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
}
