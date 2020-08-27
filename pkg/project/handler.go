package project

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset := (page - 1) * limit
	var project Project
	projectService := NewService(db.SQLLite)
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

//NameList 获取项目名称列表
func NameList(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	var project Project
	projectService := NewService(db.SQLLite)
	// 100个项目应该足够多了，先这样吧！
	list, count, err := projectService.GetProjectList(0, 100, project)
	if err != nil {
		panic(err)
	}
	var nameList []string
	for _, val := range list {
		nameList = append(nameList, val.Name)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  nameList,
		"msg":   "ok",
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
	projectService := NewService(db.SQLLite)
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
	projectService := NewService(db.SQLLite)
	project.CreatedAt = time.Now()
	err := projectService.CreateProject(&project)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    project,
		"message": "ok",
	})
}

//UpdateProject
func UpdateProject(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	projectID, _ := strconv.Atoi(c.Param("id"))
	projectService := NewService(db.SQLLite)
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
	err := ioutil.WriteFile(filePath+"/ecosystem.config.js", d1, 0666) //写入文件(字节数组)
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
	projectService := NewService(db.SQLLite)
	err := projectService.DeleteProject(projectID)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
