package exec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path"

	"github.com/yuedun/zhuque/db"
	"github.com/yuedun/zhuque/pkg/project"
	"github.com/yuedun/zhuque/util"
)

// ExecCmdSync 同步执行命令
func ExecCmdSync(userCmd string) ([]byte, error) {
	var cmdOut []byte
	var cmd *exec.Cmd
	// 执行单个shell命令时, 直接运行即可
	cmd = exec.Command("bash", "-c", userCmd)
	if cmdOut, err := cmd.CombinedOutput(); err != nil {
		log.Println("输出错误：", err)
		log.Println("输出错误2：", string(cmdOut))
		return nil, fmt.Errorf("错误码：%s，错误信息：%s", err.Error(), string(cmdOut))
	}
	return cmdOut, nil
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

// SyncCode 同步代码
func SyncCode(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
	hostLen := len(deployConfig.Host)
	ch := make(chan []byte, hostLen)
	for _, host := range deployConfig.Host {
		go func(host string) {
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
		}(host)
	}
	out := <-ch
	return out, nil
}

func Scp(projectID int) string {
	projectObj := project.Project{
		ID: projectID,
	}
	projectService := project.NewService(db.SQLLite)
	projectResult, _ := projectService.GetProjectInfo(projectObj)
	var config map[string]interface{}
	err := json.Unmarshal([]byte(projectResult.Config), &config)
	if err != nil {
		log.Println("项目配置解析失败，请检查配置json是否正确1:", err)
		panic(err)
	}
	production := config["deploy"].(map[string]interface{})
	productionJSON, err := json.Marshal(production["production"])
	var deployConfig project.DeployConfig
	err = json.Unmarshal(productionJSON, &deployConfig)
	if err != nil {
		log.Println("项目配置解析失败，请检查配置json是否正确2:", err)
		panic(err)
	}
	var buffer bytes.Buffer
	// 拉代码
	if exists := util.PathExists(path.Join(util.Conf.APPDir, projectResult.Name)); exists == false {
		log.Println(">>>>>>>>>>>", exists)
		// 分支，gitrepo，
		output, err := CloneRepo(deployConfig, projectResult.Name)
		if err != nil {
			panic(err)
		}
		buffer.Write(output)
	}
	// 装依赖
	output, err := InstallDep(deployConfig, projectResult.Name)
	if err != nil {
		panic(err)
	}
	buffer.Write(output)
	// 同步代码
	output, err = SyncCode(deployConfig, projectResult.Name)
	if err != nil {
		panic(err)
	}
	buffer.Write(output)
	return string(buffer.Bytes())
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

// PostDeploy 发布 4
func PostDeploy() ([]byte, error) {
	return nil, nil
}
