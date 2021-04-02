package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yuedun/zhuque/pkg/message"
	"github.com/yuedun/zhuque/util"
)

func TestDingTalk(t *testing.T) {
	client := &http.Client{}
	bodyObj := make(map[string]interface{})
	bodyObj["msgtype"] = "text"
	bodyObj["text"] = map[string]interface{}{
		"content": "【朱雀】我就是我, 是不一样的烟火",
	}

	Conf, err := util.GetConf("../conf.yaml")
	bytestr, _ := json.Marshal(&bodyObj)
	resp, err := client.Post(Conf.DingTalk,
		"application/json", bytes.NewBuffer(bytestr))
	if err != nil {

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

	}
	t.Log(string(body))
}

func TestDingTalk2(t *testing.T) {
	bodyObj := make(map[string]interface{})
	bodyObj["msgtype"] = "text"
	bodyObj["text"] = map[string]interface{}{
		"content": "【朱雀】我就是我, 是不一样的烟火",
	}
	Conf, _ := util.GetConf("../conf.yaml")
	messageService := message.NewMessage()
	res, _ := messageService.SendDingTalk(Conf.DingTalk, bodyObj)
	t.Log(res)
}

func TestEmail(t *testing.T) {
	Conf, _ := util.GetConf("../conf.yaml")
	util.Conf = Conf
	mailTo := strings.Split(Conf.MailTo, ";")
	t.Log("mailTo", mailTo)
	messageService := message.NewMessage()
	err := messageService.SendEmail("朱雀", "测试", mailTo)
	t.Log(err)
}
