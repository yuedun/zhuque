package util

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

var SocketCon socketio.Conn

// https://github.com/googollee/go-socket.io/blob/master/_examples/gin-gonic/main.go
func SocketEvent(router *gin.Engine) {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		SocketCon = s
		return nil
	})

	server.OnEvent("/", "msg", func(s socketio.Conn, msg string) string {
		log.Println(">>>>>>>>>接收到msg消息：", msg)
		s.SetContext(msg)
		s.Emit("msg", "这是服务端发送的消息："+msg) //接收到消息发送消息事件
		return "recv " + msg             // 接收到消息直接返回
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
