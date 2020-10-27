package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/yuedun/zhuque/util"
)

type Message interface {
	//SendDingTalk 发送钉钉消息
	SendDingTalk(dingTalkURL string, bodyObj interface{}) (dingRes DingTalkRes, err error)
	// SendEmail 发送邮件
	SendEmail(sendMsg string, to string) (emailRes EmailRes, err error)
}

/*
 * 使用不同的结构体实现接口
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

//先就这样吧
func (msg *message) SendEmail(sendMsg string, to string) (emailRes EmailRes, err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	w.WriteField("subject", "【朱雀】"+sendMsg)
	w.WriteField("to", to)
	strs := []string{
		sendMsg,
	}
	w.WriteField("content", strings.Join(strs, "\r\n"))
	w.Close()
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/mail/send-mail?token=%v", util.Conf.EmailService, util.Conf.Token), &b)
	if err != nil {
		log.Println(err)
		return emailRes, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println("发送请求失败：", err)
		return emailRes, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取body失败：", err)
		return emailRes, err
	}
	log.Println("发送邮件结果：", string(body))
	json.Unmarshal(body, &emailRes)
	return emailRes, nil
}
