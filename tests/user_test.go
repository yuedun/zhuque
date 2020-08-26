package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yuedun/zhuque/db"
	"github.com/yuedun/zhuque/pkg/user"
	"github.com/yuedun/zhuque/util"
)

func TestGetUser(t *testing.T) {
	userService := user.NewService(db.SQLLite)
	userObj := user.User{ID: 1}
	user, err := userService.GetUserInfo(userObj)
	if err != nil {
		t.Error(err)
	}
	t.Log(user)
}

func TestCreateUser(t *testing.T) {
	userService := user.NewService(db.SQLLite)
	newUser := new(user.User)
	newUser.Email = ""
	err := userService.CreateUser(newUser)
	if err != nil {
		t.Error(err)
	}
	t.Log(newUser)
}
func TestHttp(t *testing.T) {
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

func TestHttp2(t *testing.T) {
	bodyObj := make(map[string]interface{})
	bodyObj["msgtype"] = "text"
	bodyObj["text"] = map[string]interface{}{
		"content": "【朱雀】我就是我, 是不一样的烟火",
	}
	Conf, _ := util.GetConf("../conf.yaml")
	res, _ := util.DingTalk(Conf.DingTalk, bodyObj)
	t.Log(res)
}
