package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yuedun/zhuque/middleware"
	"github.com/yuedun/zhuque/pkg/exec"
	"github.com/yuedun/zhuque/pkg/project"
	"github.com/yuedun/zhuque/pkg/server"
	"github.com/yuedun/zhuque/pkg/task"
	"github.com/yuedun/zhuque/pkg/user"
)

/**
 * 路由注册
 */
func Register(router *gin.Engine) {
	userRouter := router.Group("/user")
	//user路由注册,可以给各个group加中间件
	userRouter.Use(middleware.Logger())
	{
		userRouter.GET("/list", user.List)
		//userRouter.POST("/login", user.Login)
		userRouter.POST("/login", middleware.Jwt().LoginHandler)
		userRouter.GET("/refresh_token", middleware.Jwt().RefreshHandler) // 刷新token
		userRouter.GET("/logout", middleware.Jwt().LogoutHandler)
		userRouter.GET("/info/:id", middleware.Auth(), user.GetUserInfo) //单独给某个路由添加中间件
		userRouter.GET("/users-by-sql/:id", user.GetUserInfoBySql)
		userRouter.POST("/create", user.CreateUser)
		userRouter.PUT("/update/:id", user.UpdateUser)
		userRouter.DELETE("/del/:id", user.DeleteUser)
	}

	servRouter := router.Group("/server")
	//user路由注册,可以给各个group加中间件
	servRouter.Use(middleware.Logger())
	{
		servRouter.GET("/list", server.List)
		servRouter.GET("/server-by-sql/:id", server.GetServerInfoBySql)
		servRouter.POST("/create", server.CreateServer)
		servRouter.PUT("/update/:id", server.UpdateServer)
		servRouter.DELETE("/del/:id", server.DeleteServer)
	}

	projectRouter := router.Group("/project")
	//user路由注册,可以给各个group加中间件
	projectRouter.Use(middleware.Logger())
	{
		projectRouter.GET("/list", project.List)
		projectRouter.GET("/server-by-sql/:id", project.GetProjectInfoBySql)
		projectRouter.POST("/create", project.CreateProject)
		projectRouter.PUT("/update/:id", project.UpdateProject)
		projectRouter.DELETE("/del/:id", project.DeleteProject)
	}

	deployRouter := router.Group("/deploy")
	//user路由注册,可以给各个group加中间件
	projectRouter.Use(middleware.Logger())
	{
		deployRouter.GET("/list", task.List)
		deployRouter.GET("/server-by-sql/:id", task.GetTaskInfoBySql)
		deployRouter.POST("/create", task.CreateTask)
		deployRouter.PUT("/update/:id", task.UpdateTask)
		deployRouter.DELETE("/del/:id", task.DeleteTask)
	}
	//发送命令路由注册
	execRouter := router.Group("/exec")
	{
		execRouter.POST("/send", exec.Send)
	}
}
