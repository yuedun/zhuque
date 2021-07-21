package tests

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/yuedun/zhuque/db"
	"github.com/yuedun/zhuque/pkg/permission"
	"github.com/yuedun/zhuque/pkg/task"
	"github.com/yuedun/zhuque/pkg/user"
	"github.com/yuedun/zhuque/util"
)

//！！！！重要作用，用于初始化数据库
func TestMain(m *testing.M) {
	fmt.Println("begin")
	dba, err := gorm.Open("sqlite3", "../../zhuque.db")
	dba.LogMode(true)
	db.DB = dba
	if err != nil {
		panic(err)
	}
	m.Run()
	fmt.Println("end")
}

func TestGetUser(t *testing.T) {
	userService := user.NewService(db.DB)
	userObj := user.User{ID: 5}
	user, err := userService.GetUserInfo(userObj)
	if err != nil {
		t.Error(err)
	}
	t.Log(user)
}
func TestMD5(t *testing.T) {
	p := util.GetMD5("test")
	t.Log(p)
}

func TestCreateUser(t *testing.T) {
	userService := user.NewService(db.DB)
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
	userService := user.NewService(db.DB)
	emails, err := userService.GetProjectUsersEmail("zhuque")
	if err != nil {
		t.Error(err)
	}
	t.Log("返回结果：", emails)
}

// 查询任务，关联查出用户
func TestTaskUsers(t *testing.T) {
	service := task.NewService(db.DB)
	obj := task.Task{
		ID: 329,
	}
	taskResult, err := service.GetTask(obj)
	if err != nil {
		t.Error(err)
	}
	t.Log("返回结果：", taskResult)
}

// 测试gorm Find
func TestFind(t *testing.T) {
	service := permission.NewService(db.DB)
	result, err := service.GetPermissionList()
	if err != nil {
		t.Error(err)
	}
	t.Log("返回结果：", result)
}
