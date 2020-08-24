package user

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/yuedun/zhuque/util"

	"github.com/gin-gonic/gin"
	"github.com/yuedun/zhuque/db"
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
	var user User
	userService := NewService(db.SQLLite)
	list, err := userService.GetUserList(user)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": list,
		"msg":  "ok",
	})
}

type loginData struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

//GetUserInfo
func GetUserInfo(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))
	username := c.Param("username")
	email := c.Param("email")
	userService := NewService(db.SQLLite)
	userObj := User{
		ID:       userID,
		UserName: username,
		Email:    email,
	}
	user, err := userService.GetUserInfo(userObj)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "ok",
	})
}

//GetUserInfoBySql
func GetUserInfoBySql(c *gin.Context) {
	userService := NewService(db.SQLLite)
	user, err := userService.GetUserInfoBySQL()
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

//CreateUser
func CreateUser(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	userService := NewService(db.SQLLite)
	user := User{}
	if err := c.ShouldBind(&user); err != nil {
		panic(err)
	}
	user.Password = util.GetMD5(user.UserName)
	user.Status = 1
	user.CreatedAt = time.Now()
	err := userService.CreateUser(&user)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "ok",
	})
}

//UpdateUser post json
func UpdateUser(c *gin.Context) {
	userService := NewService(db.SQLLite)
	var user User
	userID, _ := strconv.Atoi(c.Param("id"))
	//user.Addr = c.PostForm("addr")
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data":    nil,
			"message": "err",
		})
	} else {
		err := userService.UpdateUser(userID, &user)
		if err != nil {
			fmt.Println("err:", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"data":    user,
			"message": "ok",
		})
	}
}

//DeleteUser
func DeleteUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))
	userService := NewService(db.SQLLite)
	err := userService.DeleteUser(userID)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
