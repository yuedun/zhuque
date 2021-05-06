package project

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
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
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset := (page - 1) * limit
	projectName := c.Query("searchParams[projectName]")
	var project Project
	project.Name = projectName
	projectService := NewService(db.DB)
	list, count, err := projectService.GetProjectList(offset, limit, project)
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

//NameList 获取用户关联项目名称列表
func NameList(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	claims := jwt.ExtractClaims(c)
	log.Println("登录用户userid:", claims["user_id"])
	userID64 := claims["user_id"].(float64)
	userID := int(userID64)
	projectService := NewService(db.DB)
	list, err := projectService.GetProjectNameList(userID)
	if err != nil {
		panic(err)
	}
	// 分组数据
	nameList := make(map[string][]Project)
	for _, val := range list {
		// map中存在key在向该key添加数据，否则创建新key
		if v, ok := nameList[val.Namespace]; ok == true {
			v = append(v, val)
			nameList[val.Namespace] = v
		} else {
			nameList[val.Namespace] = []Project{val}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": nameList,
		"msg":  "ok",
	})
}

//NameList 获取用户关联项目名称列表 穿梭框
func NameListV2(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	claims := jwt.ExtractClaims(c)
	log.Println("登录用户userid:", claims["user_id"])
	userID64 := claims["user_id"].(float64)
	userID := int(userID64)
	projectService := NewService(db.DB)
	list, err := projectService.GetProjectNameList(userID)
	if err != nil {
		panic(err)
	}
	log.Println(list)
	// 穿梭框数据
	var nameList []map[string]string
	for _, val := range list {
		m := make(map[string]string)
		m["title"] = "[" + val.Namespace + "]" + val.Name
		m["value"] = val.Name
		nameList = append(nameList, m)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": nameList,
		"msg":  "ok",
	})
}

//NameListAll 获取所有项目名称列表
func NameListAll(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	projectService := NewService(db.DB)
	list, err := projectService.GetAllProjectNameList()
	if err != nil {
		panic(err)
	}
	// 分组数据
	nameList := make(map[string][]Project)
	for _, val := range list {
		// map中存在key在向该key添加数据，否则创建新key
		if v, ok := nameList[val.Namespace]; ok == true {
			v = append(v, val)
			nameList[val.Namespace] = v
		} else {
			nameList[val.Namespace] = []Project{val}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": nameList,
		"msg":  "ok",
	})
}

//GetProjectInfo
func GetProjectInfo(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	projectID, _ := strconv.Atoi(c.Param("id"))
	projectService := NewService(db.DB)
	projectObj := Project{
		ID: projectID,
	}
	project, err := projectService.GetProjectInfo(projectObj)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    project,
		"message": "ok",
	})
}

//CreateProject
func CreateProject(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	project := Project{}
	if err := c.ShouldBind(&project); err != nil {
		panic(err)
	}
	projectService := NewService(db.DB)
	project.CreatedAt = time.Now()
	err := projectService.CreateProject(&project)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    nil,
		"message": "ok",
	})
}

// UpdateProject 修改项目信息
func UpdateProject(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	projectID, _ := strconv.Atoi(c.Param("id"))
	projectService := NewService(db.DB)
	var project Project
	if err := c.ShouldBind(&project); err != nil {
		panic(err)
	}
	// 1.创建项目目录，2.进入目录，3.写ecosystem.config.js文件
	filePath := "./projects/" + project.Name
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		log.Println(err)
		panic(err)
	}

	var d1 = []byte(project.Config)
	fileName := "ecosystem.config.js"
	if project.DeployType == "scp" {
		fileName = "ecosystem.json"
	}
	err := ioutil.WriteFile(filePath+"/"+fileName, d1, 0666) //写入文件(字节数组)
	if err != nil {
		panic(err)
	}
	err = projectService.UpdateProject(projectID, &project)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    project,
		"message": "ok",
	})
}

//DeleteProject
func DeleteProject(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	projectID, _ := strconv.Atoi(c.Param("id"))
	projectService := NewService(db.DB)
	err := projectService.DeleteProject(projectID)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
