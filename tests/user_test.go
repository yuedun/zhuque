package tests

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/yuedun/zhuque/db"
	"github.com/yuedun/zhuque/pkg/user"
)

//！！！！重要作用，用于初始化数据库
func TestMain(m *testing.M) {
	fmt.Println("begin")
	dba, err := gorm.Open("sqlite3", "../../zhuque.db")
	dba.LogMode(true)
	db.SQLLite = dba
	if err != nil {
		panic(err)
	}
	m.Run()
	fmt.Println("end")
}

func TestGetUser(t *testing.T) {
	userService := user.NewService(db.SQLLite)
	userObj := user.User{ID: 5}
	user, err := userService.GetUserInfo(userObj)
	if err != nil {
		t.Error(err)
	}
	t.Log(user)
}

func TestCreateUser(t *testing.T) {
	userService := user.NewService(db.SQLLite)
	newUser := new(user.User)
	newUser.Email = ""
	err := userService.CreateUser(newUser)
	if err != nil {
		t.Error(err)
	}
	t.Log(newUser)
}

//查询项目关联的用户
func TestProjectUsers(t *testing.T) {
	userService := user.NewService(db.SQLLite)
	emails, err := userService.GetProjectUsersEmail("zhuque")
	if err != nil {
		t.Error(err)
	}
	t.Log("返回结果：", emails)
}
