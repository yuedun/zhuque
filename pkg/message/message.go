package message

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/yuedun/zhuque/util"
	"gopkg.in/gomail.v2"
)

type Message interface {
	//SendDingTalk 发送钉钉消息
	SendDingTalk(dingTalkURL string, bodyObj interface{}) (dingRes DingTalkRes, err error)
	// SendEmail 发送邮件
	SendEmail(subject, body string, to []string) (err error)
}

/*
 * 使用不同的结构体实现接口，如果有其他类型的邮件服务，增加结构体实现，而不是修改现有实现
 */
type message struct{}

func NewMessage() Message {
	return &message{}
}

// 钉钉消息返回
type DingTalkRes struct {
	Errcode int
	Errmsg  string
}

// 邮件返回
type EmailRes struct {
	Code int
	Msg  string
	Data string
}

/**
 * DingTalk 发送钉钉消息
 */
func (msg *message) SendDingTalk(dingTalkURL string, bodyObj interface{}) (dingRes DingTalkRes, err error) {
	client := &http.Client{}
	bytestr, _ := json.Marshal(&bodyObj)
	resp, err := client.Post(dingTalkURL,
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

// 发送邮件
func (msg *message) SendEmail(subject, body string, to []string) (err error) {
	m := gomail.NewMessage()
	// 这种方式可以添加别名，即 nickname， 也可以直接用<code>m.SetHeader("From", MAIL_USER)</code>
	nickname := "zhuque"
	// nickname := "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte("标题")) + "?="//解决中文乱码
	m.SetHeader("From", nickname+"<"+util.Conf.MailUser+">")
	// 发送给多个用户
	m.SetHeader("To", to...)
	// 设置邮件主题
	m.SetHeader("Subject", subject)
	// 设置邮件正文
	m.SetBody("text/html", body)
	d := gomail.NewDialer(util.Conf.MailHost, util.Conf.MailPort, util.Conf.MailUser, util.Conf.MailPWD)
	// 发送邮件
	err = d.DialAndSend(m)
	log.Println("发送邮件结果", err)
	return err
}
