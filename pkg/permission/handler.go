package permission

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
	userId, _ := strconv.Atoi(c.Query("userId"))
	// page, _ := strconv.Atoi(c.Query("page"))
	// limit, _ := strconv.Atoi(c.Query("limit"))
	// offset := (page - 1) * limit
	// var permission Permission
	serverService := NewService(db.SQLLite)
	list, count, err := serverService.GetPermissionList(userId)
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

//GetPermissionInfo
func GetPermissionInfo(c *gin.Context) {
	permissionID, _ := strconv.Atoi(c.Param("id"))
	// name := c.Param("name")
	permissionService := NewService(db.SQLLite)
	permissionObj := Permission{
		ID: permissionID,
	}
	permission, err := permissionService.GetPermissionInfo(permissionObj)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    permission,
		"message": "ok",
	})
}

//CreatePermission
func CreatePermission(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	permissionService := NewService(db.SQLLite)
	permission := Permission{}
	if err := c.ShouldBind(&permission); err != nil {
		panic(err)
	}
	permission.CreatedAt = time.Now()
	err := permissionService.CreatePermission(&permission)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    permission,
		"message": "ok",
	})
}

//UpdatePermission post json
func UpdatePermission(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	permissionID, _ := strconv.Atoi(c.Param("id"))
	permissionService := NewService(db.SQLLite)
	var permission Permission
	if err := c.ShouldBind(&permission); err != nil {
		panic(err)
	}
	permission.UpdatedAt = time.Now()
	err := permissionService.UpdatePermission(permissionID, &permission)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    permission,
		"message": "ok",
	})
}

//DeletePermission
func DeletePermission(c *gin.Context) {
	permissionID, _ := strconv.Atoi(c.Param("id"))
	permissionService := NewService(db.SQLLite)
	err := permissionService.DeletePermission(permissionID)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

//GetByUserID
func GetByRole(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("roleid"))
	permissionService := NewService(db.SQLLite)
	list, err := permissionService.GetByRole(roleID)
	if err != nil {
		fmt.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    list,
	})
}
