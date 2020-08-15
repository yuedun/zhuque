package tests

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yuedun/zhuque/db"
	"github.com/yuedun/zhuque/pkg/user"
)

func TestGetUser(t *testing.T) {
	userService := user.NewService(db.SQLLite)
	userObj := user.User{Id: 1}
	user, err := userService.GetUserInfo(userObj)
	if err != nil {
		t.Error(err)
	}
	t.Log(user)
}

func TestCreateUser(t *testing.T) {
	userService := user.NewService(db.SQLLite)
	newUser := new(user.User)
	newUser.Mobile = "17864345978"
	err := userService.CreateUser(newUser)
	if err != nil {
		t.Error(err)
	}
	t.Log(newUser)
}
