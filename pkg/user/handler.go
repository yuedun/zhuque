package user

import (
	"errors"
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

// List 用户列表
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
							"title":  "角色管理",
							"href":   "page/role.html",
							"icon":   "fa fa-user-circle-o",
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

// CreateUserProject 创建用户项目关系
func CreateUserProject(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	claims := jwt.ExtractClaims(c)
	createUserID64 := claims["user_id"].(float64)
	createUserID := int(createUserID64)
	log.Println("登录用户userid:", createUserID)

	userID64, _ := strconv.Atoi(c.PostForm("userID"))
	projectID64, _ := strconv.Atoi(c.PostForm("projectID"))
	log.Println(userID64, projectID64)

	userService := NewService(db.SQLLite)
	userProject := UserProject{}
	userProject.UserID = userID64
	userProject.ProjectID = projectID64
	userProject.CreateUser = createUserID
	err := userService.CreateUserProject(&userProject)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    userProject,
		"message": "ok",
	})
}

// UserProjectList 用户项目关系列表
func UserProjectList(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	userID, _ := strconv.Atoi(c.Param("userID"))
	offset := (page - 1) * limit
	var user User
	user.ID = userID
	userService := NewService(db.SQLLite)
	list, count, err := userService.GetUserProjects(offset, limit, user)
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

//DeleteUserProject 删除用户项目关系
func DeleteUserProject(c *gin.Context) {
	upID, _ := strconv.Atoi(c.Param("id"))
	userService := NewService(db.SQLLite)
	err := userService.DeleteUserProject(upID)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

//ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
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
	user := new(User)
	user.ID = userID
	oldPwd := c.PostForm("old_password")
	userService := NewService(db.SQLLite)
	uresult, err := userService.GetUserInfo(*user)
	if err != nil {
		panic(err)
	}
	if util.GetMD5(oldPwd) != uresult.Password {
		panic(errors.New("原始密码不正确！"))
	}
	user.Password = util.GetMD5(c.PostForm("new_password"))
	err = userService.UpdateUser(userID, user)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
