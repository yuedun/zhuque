package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/yuedun/zhuque/db"
	"github.com/yuedun/zhuque/pkg/permission"
	"github.com/yuedun/zhuque/pkg/project"
	"github.com/yuedun/zhuque/pkg/role"
	"github.com/yuedun/zhuque/pkg/task"
	"github.com/yuedun/zhuque/pkg/user"
	"github.com/yuedun/zhuque/router"
	"github.com/yuedun/zhuque/util"
)

// 全局配置
var Conf *util.Config

func init() {
	var err error
	Conf, err = util.GetConf("./conf.yaml")
	util.Conf = Conf
	if err != nil {
		panic(err)
	}
	db.DB, err = gorm.Open(Conf.Dialects, Conf.Dbpath)
	if err != nil {
		log.Println("failed to connect database")
		panic(err)
	}
	db.DB.AutoMigrate(&user.User{})
	db.DB.AutoMigrate(&user.UserProject{})
	db.DB.AutoMigrate(&project.Project{})
	db.DB.AutoMigrate(&task.Task{})
	db.DB.AutoMigrate(&permission.Permission{})
	db.DB.AutoMigrate(&role.Role{})
	db.DB.AutoMigrate(&role.RolePermission{})
	db.DB.LogMode(true)
	//Db.SingularTable(true) // 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
	//defer Db.Close()
}

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags) //输出行号
	r := gin.Default()
	r.Static("/fe", "./fe")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/fe")
	})

	router.Register(r)
	port := Conf.Port
	if port == "" {
		port = "8090"
	}
	log.Println("环境变量env:", util.Conf.Env, "http://localhost:"+port)
	r.Run(":" + port) // listen and serve on 0.0.0.0:8090
}
