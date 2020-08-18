package server

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type (
	/**
	  面向接口开发：
	  面向接口开发的好处是要对下面的函数进行测试时，不需要依赖一个全局的mysql连接，只需要调用NewService传一个mysql连接参数即可测试
	*/
	ServerService interface {
		GetServerInfo(search Server) (server Server, err error)
		GetServerList(search Server) (list []Server, err error)
		GetServerInfoBySQL() (server Server, err error)
		CreateServer(server *Server) (err error)
		UpdateServer(serverID int, server *Server) (err error)
		DeleteServer(serverID int) (err error)
	}
)

type svrService struct {
	mysql *gorm.DB
}

/*NewService 初始化结构体*/
func NewService(mysql *gorm.DB) ServerService {
	return &svrService{
		mysql: mysql,
	}
}

func (u *svrService) GetServerInfo(search Server) (user Server, err error) {
	err = u.mysql.Where(search).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *svrService) GetServerList(search Server) (list []Server, err error) {
	err = u.mysql.Where(search).Find(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
}

func (u *svrService) GetServerInfoBySQL() (server Server, err error) {
	err = u.mysql.Raw("select * from server where id=?", server.ID).Scan(&server).Error
	if err != nil {
		return server, err
	}
	return server, nil
}

func (u *svrService) CreateServer(user *Server) (err error) {
	err = u.mysql.Create(user).Error
	fmt.Println(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *svrService) UpdateServer(userID int, user *Server) (err error) {
	err = u.mysql.Model(user).Where("id = ?", userID).Update(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *svrService) DeleteServer(userID int) (err error) {
	u.mysql.Where("id = ?", userID).Delete(Server{})
	if err != nil {
		return err
	}
	return nil
}
