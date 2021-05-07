package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	yaml "gopkg.in/yaml.v3"
)

// 保存数据变量
var Conf *Config

//profile variables
type Config struct {
	Port        string `yaml:"port"`        // 服务端口
	Dialects    string `yaml:"dialects"`    //使用的数据库类型：mysql,sqlite3
	Dbpath      string `yaml:"dbpath"`      // sqlite数据库文件位置或mysql连接地址
	Env         string `yaml:"env"`         // 执行环境
	DingTalk    string `yaml:"dingTalk"`    // 钉钉webhook
	MailHost    string `yaml:"mailHost"`    // 邮件服务器地址
	MailPort    int    `yaml:"mailPort"`    // 邮件端口
	MailUser    string `yaml:"mailUser"`    // 邮件发送账户
	MailPWD     string `yaml:"mailPWD"`     // 邮件授权密码
	MailTo      string `yaml:"mailTo"`      // 邮件发送地址
	DelayDeploy int    `yaml:"delayDeploy"` // 延时发布时间，单位秒。默认5分钟
	JWTSecret   string `yaml:"JWTSecret"`   // jwt安全密匙
	HostName    string `yaml:"hostName"`    //服务地址
	APPDir      string `yaml:"appDir"`      // 要发布的应用存储目录
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
	if c.JWTSecret == "" {
		c.JWTSecret = "JWTSecret"
	}
	if c.APPDir == "" {
		c.APPDir = "../deploy-apps"
	}
	if c.Dialects == "" {
		c.Dialects = "sqlite3"
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
	log.Println("hex解码>>>>>>>", hex.EncodeToString(Result))
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

// CreateToken 生成jwt
func CreateToken(uid, secret string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// ParseToken 解密jwt
func ParseToken(token string, secret string) (string, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	return claim.Claims.(jwt.MapClaims)["uid"].(string), nil
}

// PathExists 判断文件或文件夹是否存在
func PathExists(path string) bool {
	log.Println("判断目录是否存在：", path)
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
