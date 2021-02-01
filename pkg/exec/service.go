package exec

import (
	"fmt"
	"log"
	"os/exec"
	"path"

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

// CloneRepo 分支，gitrepo，
func CloneRepo(deployConfig project.DeployConfig, projectName string) ([]byte, error) {
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
	// 用户名@IP:远程目录
	remotePath := fmt.Sprintf("%s@%s:%s", deployConfig.User, deployConfig.Host, deployConfig.Path)
	// rsync参数，宿主机项目，目标机地址
	cmd3 := fmt.Sprintf("rsync -av %s %s %s", deployConfig.RsyncArgs, path.Join(util.Conf.APPDir, projectName), remotePath)
	log.Println("第三步：同步代码：", cmd3)
	cmdOut, err := ExecCmdSync(cmd3)
	if err != nil {
		log.Println("第三步：同步代码执行失败：", err)
		return nil, err
	}
	return cmdOut, nil
}
