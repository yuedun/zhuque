package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v3"
)

var Conf *Config

//profile variables
type Config struct {
	Port         string `yaml:"port"`     //服务端口
	Dbpath       string `yaml:"dbpath"`   //sqlite数据库文件位置
	Env          string `yaml:"env"`      //执行环境
	DingTalk     string `yaml:"dingTalk"` //钉钉webhook
	EmailService string `yaml:"emailService"`
	Token        string `yaml:"token"`
	EmailTo      string `yaml:"emailTo"`
	TestUser     string `yaml:"testUser"`
	DelayDeploy  int    `yaml:"delayDeploy"` //延时发布时间，单位秒。默认5分钟
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
	if c.DelayDeploy == 0 {
		c.DelayDeploy = 5
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
