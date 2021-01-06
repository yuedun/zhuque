package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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

var Conf *util.Config

func init() {
	var err error
	Conf, err = util.GetConf("./conf.yaml")
	util.Conf = Conf
	if err != nil {
		panic(err)
	}
	db.SQLLite, err = gorm.Open("sqlite3", Conf.Dbpath)
	if err != nil {
		log.Println("failed to connect database")
		panic(err)
	}
	db.SQLLite.AutoMigrate(&user.User{})
	db.SQLLite.AutoMigrate(&user.UserProject{})
	db.SQLLite.AutoMigrate(&project.Project{})
	db.SQLLite.AutoMigrate(&task.Task{})
	db.SQLLite.AutoMigrate(&permission.Permission{})
	db.SQLLite.AutoMigrate(&role.Role{})
	// db.SQLLite.AutoMigrate(&role.RolePermission{})
	db.SQLLite.LogMode(true)
	//Db.SingularTable(true) // 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
	//defer Db.Close()
}

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags) //输出行号
	r := gin.Default()
	// r.LoadHTMLGlob("templates/*") //加载模板
	r.Static("/fe", "./fe")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/fe")
	})

	router.Register(r)
	port := Conf.Port
	if port == "" {
		port = "8090"
	}
	log.Println("环境变量env:", util.Conf.Env)
	r.Run(":" + port) // listen and serve on 0.0.0.0:8090
}
