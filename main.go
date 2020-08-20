package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/yuedun/zhuque/db"
	"github.com/yuedun/zhuque/pkg/project"
	"github.com/yuedun/zhuque/pkg/server"
	"github.com/yuedun/zhuque/pkg/task"
	"github.com/yuedun/zhuque/pkg/user"
	"github.com/yuedun/zhuque/router"
	"github.com/yuedun/zhuque/util"
)

func init() {
	var err error
	var conf util.Conf
	c, err := conf.GetConf("./conf.yaml")
	if err != nil {
		panic(err)
	}
	db.SQLLite, err = gorm.Open("sqlite3", c.Dbpath)
	if err != nil {
		log.Println("failed to connect database")
		panic(err)
	}
	db.SQLLite.AutoMigrate(&user.User{})
	db.SQLLite.AutoMigrate(&server.Server{})
	db.SQLLite.AutoMigrate(&project.Project{})
	db.SQLLite.AutoMigrate(&task.Task{})
	db.SQLLite.LogMode(true)
	//Db.SingularTable(true) // 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
	//defer Db.Close()
}

func main() {
	r := gin.Default()
	//r.Use(middleware.Logger())//全局中间件
	//r.LoadHTMLGlob("templates/*") //加载模板
	r.Static("/fe", "./fe")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/fe")
	})

	router.Register(r)
	r.Run(":8090") // listen and serve on 0.0.0.0:8090
}