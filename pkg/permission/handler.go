package permission

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yuedun/zhuque/db"
)

//List 获取所有权限列表
func List(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	serverService := NewService(db.DB)
	list, err := serverService.GetPermissionList()
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": list,
		"msg":  "ok",
	})
}

//GetPermissionInfo
func GetPermissionInfo(c *gin.Context) {
	permissionID, _ := strconv.Atoi(c.Param("id"))
	// name := c.Param("name")
	permissionService := NewService(db.DB)
	permissionObj := Permission{
		ID: permissionID,
	}
	permission, err := permissionService.GetPermissionInfo(permissionObj)
	if err != nil {
		log.Println("err:", err)
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
	permissionService := NewService(db.DB)
	permission := Permission{}
	if err := c.ShouldBind(&permission); err != nil {
		panic(err)
	}
	err := permissionService.CreatePermission(&permission)
	if err != nil {
		log.Println("err:", err)
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
	permissionService := NewService(db.DB)
	var permission Permission
	if err := c.ShouldBind(&permission); err != nil {
		panic(err)
	}
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
	permissionService := NewService(db.DB)
	err := permissionService.DeletePermission(permissionID)
	if err != nil {
		log.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

//RolePermissions 角色管理-分配权限
func RolePermissions(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("roleid"))
	permissionService := NewService(db.DB)
	//所有权限
	allPermissionList, err := permissionService.GetPermissionListForRole()
	//角色已分配权限
	rolePermissionList, err := permissionService.GetByRole(roleID)
	if err != nil {
		log.Println("err:", err)
	}
	for _, per := range allPermissionList {
		// permissionTree.Field = per.MenuURL
		per.Spread = true
		for _, rp := range rolePermissionList {
			if per.ID == rp.ID {
				per.Checked = true
				// break
			}
			for _, child := range per.Children {
				if rp.ID == child.ID {
					child.Checked = true
					per.Checked = false
					//layui tree插件父节点选中会忽略子节点是否选中，默认全选中。但子节点选中父节点也会选中，所以子节点选中时设置父节点为不选中，以免子节点全部选中。
				}
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    allPermissionList,
	})
}
