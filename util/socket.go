package util

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

var SocketConMap = make(map[string]socketio.Conn)

// https://github.com/googollee/go-socket.io/blob/master/_examples/gin-gonic/main.go
func SocketEvent(router *gin.Engine) {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		userID := NewCookie(s.RemoteHeader().Get("Cookie")).getCookieByName("userID")
		fmt.Println("connected:", s.ID(), userID)
		SocketConMap[userID] = s
		return nil
	})

	server.OnEvent("/", "msg", func(s socketio.Conn, msg string) string {
		log.Println(">>>>>>>>>接收到msg消息：", msg)
		s.SetContext(msg)
		time := time.Now().Format("15:04:05.000")
		s.Emit("msg", SocketData{Time: time, Msg: "socket连接服务器成功！"}) //接收到消息发送消息事件
		return "recv " + msg                                         // 接收到消息直接返回
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed:", reason)
	})

	go server.Serve()
	// defer server.Close()//panic: send on closed channel
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))

}

type myCookie struct {
	cookies []*http.Cookie
}

// 将字符串cookie转换为数组cookie
func NewCookie(rawCookies string) myCookie {
	header := http.Header{}
	header.Add("Cookie", rawCookies)
	request := http.Request{Header: header}

	return myCookie{cookies: request.Cookies()}
}

// 获取指定key的cookie值
func (c myCookie) getCookieByName(key string) string {
	val := ""
	for _, v := range c.cookies {
		if v.Name == key {
			val = v.Value
			break
		}
	}
	return val
}

type SocketData struct {
	Time   string
	Msg    string
	Status int
}

func SendSocketMsg(uid, msg string) {
	if v, ok := SocketConMap[uid]; ok {
		time := time.Now().Format("15:04:05.000")
		v.Emit("msg", SocketData{Time: time, Msg: msg})
	}
}
