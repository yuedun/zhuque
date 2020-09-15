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

// Register 路由注册
func Register(router *gin.Engine) {
	userRouter := router.Group("/user")
	//user路由注册,可以给各个group加中间件
	userRouter.Use(middleware.Logger())
	{
		userRouter.POST("/login", middleware.Jwt().LoginHandler)
		userRouter.GET("/refresh_token", middleware.Jwt().RefreshHandler) // 刷新token
		userRouter.GET("/logout", middleware.Jwt().LogoutHandler)
		userRouter.GET("/info/:id", middleware.Jwt().MiddlewareFunc(), user.GetUserInfo) //单独给某个路由添加中间件
		userRouter.GET("/list", user.List)
		userRouter.POST("/create", user.CreateUser)
		userRouter.PUT("/update/:id", user.UpdateUser)
		userRouter.DELETE("/del/:id", user.DeleteUser)
		userRouter.GET("/init", middleware.Jwt().MiddlewareFunc(), user.Init)
		userRouter.GET("/user-projects/:userID", middleware.Jwt().MiddlewareFunc(), user.UserProjectList)
		userRouter.POST("/create-user-project", middleware.Jwt().MiddlewareFunc(), user.CreateUserProject)
	}

	servRouter := router.Group("/server")
	//user路由注册,可以给各个group加中间件
	servRouter.Use(middleware.Logger())
	servRouter.Use(middleware.Jwt().MiddlewareFunc())
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
	projectRouter.Use(middleware.Jwt().MiddlewareFunc())
	{
		projectRouter.GET("/list", project.List)
		projectRouter.GET("/name-list", project.NameList)
		projectRouter.GET("/name-list-v2", project.NameListV2)
		projectRouter.POST("/create", project.CreateProject)
		projectRouter.GET("/get-by-id/:id", project.GetProjectInfo)
		projectRouter.PUT("/update/:id", project.UpdateProject)
		projectRouter.POST("/del/:id", project.DeleteProject)
	}

	deployRouter := router.Group("/deploy")
	//user路由注册,可以给各个group加中间件
	projectRouter.Use(middleware.Logger())
	projectRouter.Use(middleware.Jwt().MiddlewareFunc())
	{
		deployRouter.GET("/list", task.List)
		deployRouter.GET("/wait-list", task.WaitList)
		deployRouter.GET("/server-by-sql/:id", task.GetTaskInfoBySql)
		deployRouter.POST("/create", task.CreateTask)
		deployRouter.PUT("/update/:id", task.UpdateTask)
		deployRouter.DELETE("/del/:id", task.DeleteTask)
		deployRouter.POST("/approve", task.Approve)
	}
	//发送命令路由注册
	execRouter := router.Group("/exec")
	execRouter.Use(middleware.Jwt().MiddlewareFunc())
	{
		execRouter.POST("/send", exec.Send)
		execRouter.POST("/server", exec.Server)
		execRouter.POST("/server-v2", exec.ServerV2)
		execRouter.POST("/release/:id", exec.Release)
		execRouter.POST("/release-v2/:id", exec.ReleaseV2)
	}
}
