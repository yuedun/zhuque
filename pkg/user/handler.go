package user

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
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
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset := (page - 1) * limit
	username := c.Query("searchParams[username]")
	email := c.Query("searchParams[email]")
	var user User
	user.UserName = username
	user.Email = email
	userService := NewService(db.SQLLite)
	list, count, err := userService.GetUserList(offset, limit, user)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  list,
		"msg":   "ok",
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

// Init 初始化用户界面
func Init(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	claims := jwt.ExtractClaims(c)
	userID64 := claims["user_id"].(float64)
	userID := int(userID64)
	log.Println("登录用户userid:", userID)
	userService := NewService(db.SQLLite)
	user := User{ID: userID}
	userObj, err := userService.GetUserInfo(user)
	if err != nil {
		fmt.Println("err:", err)
		panic(err)
	}
	c.JSON(http.StatusOK,
		map[string]interface{}{
			"homeInfo": map[string]string{
				"title": "首页",
				"href":  "page/welcome.html?t=1",
			},
			"logoInfo": map[string]string{
				"title": "朱雀",
				"image": "images/logo.png",
				"href":  "",
			},
			"userInfo": map[string]interface{}{
				"username": userObj.UserName,
				"userID":   userObj.ID,
				"role":     userObj.Role,
			},
			"menuInfo": []map[string]interface{}{
				{
					"title":  "常规管理",
					"icon":   "fa fa-address-book",
					"href":   "",
					"target": "_self",
					"child": []map[string]interface{}{
						{
							"title":  "快捷发布",
							"href":   "page/quick-release.html",
							"icon":   "fa fa-bolt",
							"target": "_self",
						},
						{
							"title":  "快捷发布-多项目",
							"href":   "page/quick-release-v2.html",
							"icon":   "fa fa-bolt",
							"target": "_self",
						},
						{
							"title":  "发布记录",
							"href":   "page/deploy.html",
							"icon":   "fa fa-tasks",
							"target": "_self",
						},
						{
							"title":  "命令部署",
							"href":   "page/task.html",
							"icon":   "fa fa-adjust",
							"target": "_self",
						},
						{
							"title":  "项目管理",
							"href":   "page/projects.html",
							"icon":   "fa fa-navicon",
							"target": "_self",
						},
						{
							"title":  "用户管理",
							"href":   "page/users.html",
							"icon":   "fa fa-users",
							"target": "_self",
						},
						{
							"title":  "菜单管理",
							"href":   "page/menu.html",
							"icon":   "fa fa-list-alt",
							"target": "_self",
						},
					},
				},
			},
		})
}
