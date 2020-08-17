package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/yuedun/zhuque/router"
)

//func init() {
//	var err error
//	db.SQLLite, err = gorm.Open("sqlite3", "../zhuque.db")
//	if err != nil {
//		panic("failed to connect database")
//	}
//	db.SQLLite.AutoMigrate(&user.User{})
//	db.SQLLite.LogMode(true)
//	//Db.SingularTable(true) // 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
//	//defer Db.Close()
//}

func main() {
	r := gin.Default()
	//r.Use(middleware.Logger())//全局中间件
	//r.LoadHTMLGlob("templates/*") //加载模板
	r.Static("/fe","./fe")
	r.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	router.Register(r)
	r.Run(":8090") // listen and serve on 0.0.0.0:8090
}
