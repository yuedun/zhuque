package message

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Message interface {
	//SendDingTalk 发送钉钉消息
	SendDingTalk(dingTalkURL string, bodyObj interface{}) (dingRes DingTalkRes, err error)
	// SendEmail 发送邮件
	SendEmail(dingTalkURL string, bodyObj interface{}) (dingRes DingTalkRes, err error)
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

/**
 * DingTalk 发送钉钉消息
 */
func (msg *message) SendDingTalk(dingTalkUrl string, bodyObj interface{}) (dingRes DingTalkRes, err error) {
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

func (msg *message) SendEmail(dingTalkURL string, bodyObj interface{}) (dingRes DingTalkRes, err error) {
	return dingRes, nil
}
