package tests

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCmd(t *testing.T) {
	var whoami []byte
	var err error
	var cmd *exec.Cmd
	// 执行单个shell命令时, 直接运行即可
	cmd = exec.Command("whoami")
	if whoami, err = cmd.Output(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// 默认输出有一个换行
	fmt.Println(string(whoami))
	// 指定参数后过滤换行符
	fmt.Println(strings.Trim(string(whoami), "\n"))

	fmt.Println("====")
}

func TestInput(t *testing.T) {
	cmd := exec.Command("tr", "a-z", "A-Z")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("in all caps: %q\n", out.String())
}

func TestCat(t *testing.T) {
	cmd := exec.Command("cat")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", out)
}

func TestFile(t *testing.T) {
	_, err := os.Stat("../../apps/vx-1")
	// t.Log(f.IsDir())
	t.Log(err)
	if err != nil {
		if os.IsExist(err) {
			t.Log("存在")
			return
		}
		t.Log("不存在")
	}
}

func TestSubstr(t *testing.T) {
	str := "删掉了改好aaa中文b中文二bb;"
	po := strings.LastIndexAny(str, ";")
	t.Log(po)
	ss := str[:po]
	t.Log(po, ss)
}

func TestChan(t *testing.T) {
	ch := make(chan []byte, 3)
	// errch := make(chan error, 3)
	hosts := []string{"a", "b", "c"}
	for _, host := range hosts {
		go func(host string, ch chan []byte) {
			t.Log("host:", host)
			ch <- []byte(host)
		}(host, ch)
	}

	i := 0
	for {
		select {
		case out := <-ch:
			i++
			t.Log(string(out))
		}
		if i == 3 {
			goto L
		}
	}
L:
	t.Log("out")
}
