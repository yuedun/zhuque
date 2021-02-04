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
	// 拉代码
	if exists := util.PathExists(path.Join(util.Conf.APPDir, projectResult.Name)); exists == false {
		// 分支，gitrepo，
		output, err := CloneRepo(deployConfig, projectResult.Name)
		if err != nil {
			return "", err
		}
		buffer.Write(output)
	}
	// 装依赖
	output, err = InstallDep(deployConfig, projectResult.Name)
	if err != nil {
		return "", err
	}
	buffer.Write(output)
	// 同步代码到远程服务器
	// output, err = SyncCode(deployConfig, projectResult.Name)
	// if err != nil {
	// 	return "", err
	// }
	// buffer.Write(output)

	// 同步代码到远程应用服务器后执行命令，如重启
	if deployConfig.PostDeploy != "" {
		output, err = PostDeploy(deployConfig)
		if err != nil {
			return "", err
		}
		buffer.Write(output)
	}
	return string(buffer.Bytes()), nil
}

// ExecCmdSync 同步执行命令
func ExecCmdSync(userCmd string) ([]byte, error) {
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
	log.Println(userCmd, "执行结果：", string(stdoutStderr))
	return stdoutStderr, nil
}

// CloneRepo clone代码
func CloneRepo(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
	//分支，gitrepo，目录名
	cmd1 := fmt.Sprintf("git clone -b %s %s %s", deployConfig.Ref, deployConfig.Repo, path.Join(util.Conf.APPDir, projectName))
	log.Println("第一步：检出代码：", cmd1)
	cmdOut, err := ExecCmdSync(cmd1)
	if err != nil {
		log.Println("第一步：检出代码执行失败：", err)
		return nil, err
	}
	return cmdOut, nil
}

// InstallDep 安装依赖
func InstallDep(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
	cmd2 := "cd " + path.Join(util.Conf.APPDir, projectName) + " && npm i"
	log.Println("第二步：安装依赖：", cmd2)
	cmdOut, err := ExecCmdSync(cmd2)
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
			cmdput, err := ExecCmdSync(cmd3)
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
		if i == 3 {
			goto L
		}
	}
L:
	log.Println("out")

	return buffer.Bytes(), outErr
}

// PostDeploy 5 远程应用服务器上发布 用户自定义命令：比如重启服务
// 利用ssh在远程服务器上执行
func PostDeploy(deployConfig project.DeployConfig) ([]byte, error) {
	hostLen := len(deployConfig.Host)
	ch := make(chan []byte, hostLen)
	// errch := make(chan error) //不需要缓冲，只要接收到一个错误就退出
	for _, host := range deployConfig.Host {
		go func(host string, ch chan []byte) {
			// 用户名@IP 命令
			ssh := fmt.Sprintf("ssh %s@%s \"%s\"", deployConfig.User, host, deployConfig.PostDeploy)
			cmdput, err := ExecCmdSync(ssh)
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
	var outErr error
	for {
		select {
		case out := <-ch:
			i++
			buffer.Write(out)
			log.Println(string(out))
		}
		if i == 3 {
			goto L
		}
	}
L:
	log.Println("out")

	return buffer.Bytes(), outErr
}
