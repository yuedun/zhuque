package exec

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path"

	"github.com/yuedun/zhuque/db"
	"github.com/yuedun/zhuque/pkg/project"
	"github.com/yuedun/zhuque/util"
)

// DeployControl 发布流程控制
func DeployControl(projectID int) (string, error) {
	projectObj := project.Project{
		ID: projectID,
	}
	projectService := project.NewService(db.SQLLite)
	projectResult, _ := projectService.GetProjectInfo(projectObj)
	var config map[string]interface{}
	err := json.Unmarshal([]byte(projectResult.Config), &config)
	if err != nil {
		log.Println("项目配置解析失败，请检查配置json是否正确1:", err)
		return "", err
	}
	production := config["deploy"].(map[string]interface{})
	productionJSON, err := json.Marshal(production["production"])
	var deployConfig project.DeployConfig
	err = json.Unmarshal(productionJSON, &deployConfig)
	if err != nil {
		log.Println("项目配置解析失败，请检查配置json是否正确2:", err)
		return "", err
	}
	var buffer bytes.Buffer
	var output []byte

	// 克隆代码
	if exists := util.PathExists(path.Join(util.Conf.APPDir, projectResult.Name)); exists == false {
		// 分支，gitrepo，
		output, err = CloneRepo(deployConfig, projectResult.Name)
		if err != nil {
			return "", err
		}
		buffer.Write(output)
	} else {
		buffer.Write([]byte("项目已存在，跳过克隆代码。\n"))
	}

	// 拉新代码
	output, err = GitPull(deployConfig, projectResult.Name)
	if err != nil {
		return "", err
	}
	buffer.Write(output)

	// 装依赖
	output, err = InstallDep(deployConfig, projectResult.Name)
	if err != nil {
		return "", err
	}
	buffer.Write(output)

	// 编译
	if deployConfig.Build != "" {
		output, err = Build(deployConfig, projectResult.Name)
		if err != nil {
			return "", err
		}
		buffer.Write(output)
	}

	// 同步代码到远程服务器
	output, err = SyncCode(deployConfig, projectResult.Name)
	if err != nil {
		return "", err
	}
	buffer.Write(output)

	// 同步代码到远程应用服务器后执行命令，如重启
	if deployConfig.PostDeploy != "" {
		output, err = PostDeploy(deployConfig, projectResult.Name)
		if err != nil {
			log.Println("远程命令执行异常：", err)
			return "", err
		}
		buffer.Write(output)
	}
	return string(buffer.Bytes()), nil
}

// CmdSync 同步执行命令
func CmdSync(userCmd string) ([]byte, error) {
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
func CloneRepo(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
	//分支，gitrepo，目录名
	cmd1 := fmt.Sprintf("git clone -b %s %s %s", deployConfig.Ref, deployConfig.Repo, path.Join(util.Conf.APPDir, projectName))
	log.Println("第一步：检出代码：", cmd1)
	cmdOut, err := CmdSync(cmd1)
	if err != nil {
		log.Println("第一步：检出代码执行失败：", err)
		return nil, err
	}
	return cmdOut, nil
}

// InstallDep 安装依赖
func InstallDep(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
	cmd := fmt.Sprintf("cd %s ; npm i", path.Join(util.Conf.APPDir, projectName))
	log.Println("第二步：安装依赖：", cmd)
	cmdOut, err := CmdSync(cmd)
	if err != nil {
		log.Println("第二步：安装依赖执行失败：", err)
		return nil, err
	}
	return cmdOut, nil
}

// PreSetup 1
func PreSetup() ([]byte, error) {
	return nil, nil
}

// PostSetup 2
func PostSetup() ([]byte, error) {
	return nil, nil
}

// PreDeployLocal 3
func PreDeployLocal() ([]byte, error) {
	return nil, nil
}

// GitPull 3
func GitPull(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
	gitpull := fmt.Sprintf("cd %s; git pull origin %s; git log --oneline -1", path.Join(util.Conf.APPDir, projectName), deployConfig.Ref)
	cmdOut, err := CmdSync(gitpull)
	if err != nil {
		log.Println("拉取代码失败：", err)
		return nil, err
	}
	return cmdOut, nil
}

// Build 3
func Build(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
	build := fmt.Sprintf("cd %s; %s", path.Join(util.Conf.APPDir, projectName), deployConfig.Build)
	cmdOut, err := CmdSync(build)
	if err != nil {
		log.Println("拉取代码失败：", err)
		return nil, err
	}
	return cmdOut, nil
}

// SyncCode 4 同步代码
func SyncCode(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
	hostLen := len(deployConfig.Host)
	ch := make(chan []byte, hostLen)
	for _, host := range deployConfig.Host {
		go func(host string, ch chan []byte) {
			// 用户名@IP:远程目录
			remotePath := fmt.Sprintf("%s@%s:%s", deployConfig.User, host, deployConfig.Path)
			// rsync参数，宿主机项目，目标机地址
			cmd3 := fmt.Sprintf("rsync -av %s %s %s", deployConfig.RsyncArgs, path.Join(util.Conf.APPDir, projectName), remotePath)
			log.Println("第三步：同步代码：", cmd3)
			cmdput, err := CmdSync(cmd3)
			if err != nil {
				log.Println("第三步：同步代码执行失败：", err)
				ch <- []byte(err.Error())
				return
			}
			ch <- cmdput
		}(host, ch)
	}
	i := 0
	var buffer bytes.Buffer
	var outErr error
	for {
		select {
		case out := <-ch:
			i++
			buffer.Write(out)
			log.Println(string(out))
		}
		if i == hostLen {
			goto L
		}
	}
L:
	log.Println("out")

	return buffer.Bytes(), outErr
}

// PostDeploy 5 远程应用服务器上发布 用户自定义命令：比如重启服务
// 利用ssh在远程服务器上执行
func PostDeploy(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
	hostLen := len(deployConfig.Host)
	ch := make(chan []byte, hostLen)
	// errch := make(chan error) //不需要缓冲，只要接收到一个错误就退出
	for _, host := range deployConfig.Host {
		go func(host string, ch chan []byte) {
			// 用户名@IP 命令
			ssh := fmt.Sprintf("ssh %s@%s \"cd %s; %s\"", deployConfig.User, host, path.Join(deployConfig.Path, projectName), deployConfig.PostDeploy)
			cmdput, err := CmdSync(ssh)
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
	log.Println("out")

	return buffer.Bytes(), nil
}
