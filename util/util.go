package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	yaml "gopkg.in/yaml.v3"
)

var Conf *Config

//profile variables
type Config struct {
	Port     string `yaml:"port"`     //服务端口
	Dbpath   string `yaml:"dbpath"`   //sqlite数据库文件位置
	Env      string `yaml:"env"`      //执行环境
	DingTalk string `yaml:"dingTalk"` //钉钉webhook
}

func GetConf(filename string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var c = new(Config)
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return c, nil
}

/**
 * md5加密
 */
func GetMD5(password string) string {
	Md5Inst := md5.New()
	Md5Inst.Write([]byte(password))
	Result := Md5Inst.Sum(nil)
	// 以下两种输出结果一样
	fmt.Println("格式化>>>>>>>%x\n", Result)
	fmt.Println("hex解码>>>>>>>", hex.EncodeToString(Result), "\n")
	return fmt.Sprintf("%x", Result)
}

/**
 * 生成密码
 */
func GeneratePassword(mobile string) string {
	b := []byte(mobile)
	p := b[7:]
	password := "hello" + string(p)
	return GetMD5(password)
}

/**
 * DingTalk 发送钉钉消息
 */
func SendDingTalk(dingTalkUrl string, bodyObj interface{}) (dingRes DingTalkRes, err error) {
	client := &http.Client{}
	bytestr, _ := json.Marshal(&bodyObj)
	resp, err := client.Post(dingTalkUrl,
		"application/json", bytes.NewBuffer(bytestr))
	if err != nil {
		return dingRes, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return dingRes, err
	}
	json.Unmarshal(body, &dingRes)
	return dingRes, nil
}

// 钉钉消息返回
type DingTalkRes struct {
	Errcode int
	Errmsg  string
}
