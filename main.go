package main

import (
	"log"
	"net/http"
	"time"

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
		log.Println("failed to open conf.yaml")
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
	initData()
}

func initData() {
	users := make([]user.User, 0)
	if err := db.DB.Model(&user.User{}).Find(&users).Error; err != nil {
		log.Println(err)
	} else {
		log.Println(">>>>>>>>", users)
		if len(users) == 0 {
			userService := user.NewService(db.DB)
			user := user.User{}
			user.UserName = "test"
			user.Password = util.GetMD5(user.UserName)
			user.Status = 1
			user.CreatedAt = time.Now()
			userService.CreateUser(&user)
		}
	}
}

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags) //输出行号
	r := gin.Default()
	r.Static("/fe", "./fe")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/fe")
	})

	util.SocketEvent(r)
	router.Register(r)
	port := Conf.Port
	if port == "" {
		port = "8090"
	}
	log.Println("环境变量env:", util.Conf.Env, "http://localhost:"+port)
	r.Run(":" + port) // listen and serve on 0.0.0.0:8090
}

// 以下代码为go:embed打包方式，将静态文件一同打包到可执行文件中
// //go:embed fe/*
// var f embed.FS

// func main() {
// 	log.SetFlags(log.Llongfile | log.LstdFlags) //输出行号
// 	r := gin.Default()
// 	templ := template.Must(template.New("").ParseFS(f, "fe/index.html"))
// 	r.SetHTMLTemplate(templ)
// 	r.StaticFS("/static", http.FS(f))
// 	r.GET("/", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "index.html", gin.H{})
// 	})

// 	util.SocketEvent(r)
// 	router.Register(r)
// 	port := Conf.Port
// 	if port == "" {
// 		port = "8090"
// 	}
// 	log.Println("环境变量env:", util.Conf.Env, "http://localhost:"+port)
// 	r.Run(":" + port) // listen and serve on 0.0.0.0:8090
// }
