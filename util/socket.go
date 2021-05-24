package util

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

var SocketCon socketio.Conn

var SocketConMap = make(map[string]socketio.Conn)

// https://github.com/googollee/go-socket.io/blob/master/_examples/gin-gonic/main.go
func SocketEvent(router *gin.Engine) {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		userID := getCookieByName("userID", getCookies(s.RemoteHeader().Get("Cookie")))
		fmt.Println("connected:", s.ID(), userID)
		// SocketCon = s
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

func getCookies(rawCookies string) []*http.Cookie {
	// rawCookies := "cookie1=value1;cookie2=value2"

	header := http.Header{}
	header.Add("Cookie", rawCookies)
	request := http.Request{Header: header}

	// fmt.Println(request.Cookies()) // [cookie1=value1 cookie2=value2]
	return request.Cookies()
}

func getCookieByName(key string, cookies []*http.Cookie) string {
	val := ""
	for _, v := range cookies {
		if v.Name == key {
			val = v.Value
			break
		}
	}
	return val
}

type SocketData struct {
	Time string
	Msg  string
}
