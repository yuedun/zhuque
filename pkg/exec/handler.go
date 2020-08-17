package exec

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"os/exec"
	"strings"
)

//发送命令到服务器
func Exec(c *gin.Context) {
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
	log.Println(">>>>>>>>>>", userCmd)
	var cmdOut []byte
	var err error
	var cmd *exec.Cmd
	// 执行单个shell命令时, 直接运行即可
	cmd = exec.Command("bash", "-c", userCmd)
	if cmdOut, err = cmd.Output(); err != nil {
		log.Println(err)
		panic(err)
	}
	// 默认输出有一个换行
	log.Println(string(cmdOut))
	// 指定参数后过滤换行符
	log.Println(strings.Trim(string(cmdOut), "\n"))

	c.JSON(200, gin.H{
		"message": "ok",
		"data":    string(cmdOut),
	})
}
