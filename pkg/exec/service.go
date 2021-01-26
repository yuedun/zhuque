package exec

import (
	"fmt"
	"log"
	"os/exec"
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
