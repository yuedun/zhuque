package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/yuedun/zhuque/db"

	"github.com/gin-gonic/gin"
)

//List
func List(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	var server Server
	serverService := NewService(db.SQLLite)
	list, err := serverService.GetServerList(server)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data":    list,
		"msg": "ok",
	})
}

//GetServerInfo
func GetServerInfo(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))
	name := c.Param("name")
	ip := c.Param("ip")
	userService := NewService(db.SQLLite)
	userObj := Server{
		ID:       userID,
		Name: name,
		IP:   ip,
	}
	user, err := userService.GetServerInfo(userObj)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "ok",
	})
}

//GetServerInfoBySql
func GetServerInfoBySql(c *gin.Context) {
	userService := NewService(db.SQLLite)
	user, err := userService.GetServerInfoBySQL()
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

//CreateServer
func CreateServer(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	userService := NewService(db.SQLLite)
	user := Server{}
	if err := c.ShouldBind(&user); err != nil {
		panic(err)
	}
	user.CreatedAt = time.Now()
	err := userService.CreateServer(&user)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "ok",
	})
}

//UpdateServer post json
func UpdateServer(c *gin.Context) {
	userService := NewService(db.SQLLite)
	var user Server
	userID, _ := strconv.Atoi(c.Param("id"))
	//user.Addr = c.PostForm("addr")
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data":    nil,
			"message": "err",
		})
	} else {
		err := userService.UpdateServer(userID, &user)
		if err != nil {
			fmt.Println("err:", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"data":    user,
			"message": "ok",
		})
	}
}

//DeleteServer
func DeleteServer(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))
	userService := NewService(db.SQLLite)
	err := userService.DeleteServer(userID)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
