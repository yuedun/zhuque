package role

import (
	"log"
	"net/http"
	"strconv"

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
	// userId, _ := strconv.Atoi(c.Query("userId"))
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset := (page - 1) * limit
	var role Role
	serverService := NewService(db.DB)
	list, count, err := serverService.GetRoleList(offset, limit, role)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"data":  list,
		"count": count,
		"msg":   "ok",
	})
}

//GetRoleInfo
func GetRoleInfo(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("id"))
	// name := c.Param("name")
	roleService := NewService(db.DB)
	roleObj := Role{
		RoleNum: roleID,
	}
	role, err := roleService.GetRoleInfo(roleObj)
	if err != nil {
		log.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    role,
		"message": "ok",
	})
}

//CreateRole
func CreateRole(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	roleService := NewService(db.DB)
	role := Role{}
	if err := c.ShouldBind(&role); err != nil {
		panic(err)
	}
	err := roleService.CreateRole(&role)
	if err != nil {
		log.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    role,
		"message": "ok",
	})
}

//SetPermission 管理员设置用户权限
func SetPermission(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	roleService := NewService(db.DB)
	role := Role{}
	if err := c.ShouldBind(&role); err != nil {
		panic(err)
	}
	err := roleService.UpdateRole(role.ID, &role)
	if err != nil {
		log.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    role,
		"message": "ok",
	})
}

//UpdateRole post json
func UpdateRole(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.(error).Error(),
			})
		}
	}()
	roleID, _ := strconv.Atoi(c.Param("id"))
	roleService := NewService(db.DB)
	var role Role
	if err := c.ShouldBind(&role); err != nil {
		panic(err)
	}
	err := roleService.UpdateRole(roleID, &role)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    role,
		"message": "ok",
	})
}

//DeleteRole
func DeleteRole(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("id"))
	roleService := NewService(db.DB)
	err := roleService.DeleteRole(roleID)
	if err != nil {
		log.Println("err:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
