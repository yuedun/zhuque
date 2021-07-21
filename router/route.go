package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yuedun/zhuque/middleware"
	"github.com/yuedun/zhuque/pkg/exec"
	"github.com/yuedun/zhuque/pkg/permission"
	"github.com/yuedun/zhuque/pkg/project"
	"github.com/yuedun/zhuque/pkg/role"
	"github.com/yuedun/zhuque/pkg/task"
	"github.com/yuedun/zhuque/pkg/user"
)

// Register 路由注册
func Register(router *gin.Engine) {
	router.Use(middleware.Logger()) //全局中间件
	// router.Use(middleware.SetUserInfo()) //这个中间件不能加在这，因为获取用户信息的前置是已经设置了用户信息，而设置用户信息是在middleware.Jwt().MiddlewareFunc()中间件中操作的
	userRouter := router.Group("/user")
	//user路由注册,可以给各个group加中间件
	userRouter.POST("/login", middleware.Jwt().LoginHandler)
	userRouter.GET("/refresh_token", middleware.Jwt().RefreshHandler) // 刷新token
	userRouter.GET("/logout", middleware.Jwt().LogoutHandler)
	userRouter.GET("/forgot-password", user.ForgotPassword) //忘记密码，发送邮件
	userRouter.GET("/reset-password", user.RestPassword)    //根据邮件链接重置密码
	userRouter.Use(middleware.Jwt().MiddlewareFunc())
	{
		userRouter.GET("/info/:id", user.GetUserInfo) //单独给某个路由添加中间件
		userRouter.GET("/list", user.List)
		userRouter.POST("/create", user.CreateUser)
		userRouter.PUT("/update/:id", user.UpdateUser)
		userRouter.DELETE("/del/:id", user.DeleteUser)
		userRouter.GET("/init", user.Init)
		userRouter.GET("/user-projects/:userID", user.UserProjectList)
		userRouter.POST("/create-user-project", user.CreateUserProject)
		userRouter.DELETE("/user-project/del/:id", user.DeleteUserProject)
		userRouter.POST("/change-password", user.ChangePassword)
	}

	projectRouter := router.Group("/project")
	//user路由注册,可以给各个group加中间件
	projectRouter.Use(middleware.Jwt().MiddlewareFunc())
	projectRouter.Use(middleware.SetUserInfo()) // 自动设置登录用的的id
	{
		projectRouter.GET("/list", project.List)
		projectRouter.GET("/name-list", project.NameList)
		projectRouter.GET("/name-list-v2", project.NameListV2)
		projectRouter.GET("/name-list-all", project.NameListAll)
		projectRouter.POST("/create", project.CreateProject)
		projectRouter.GET("/get-by-id/:id", project.GetProjectInfo)
		projectRouter.PUT("/update/:id", project.UpdateProject)
		projectRouter.POST("/del/:id", project.DeleteProject)
	}

	deployRouter := router.Group("/deploy")
	//user路由注册,可以给各个group加中间件
	deployRouter.Use(middleware.Jwt().MiddlewareFunc())
	deployRouter.Use(middleware.SetUserInfo()) // 自动设置登录用的的id
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
	execRouter.Use(middleware.SetUserInfo()) // 自动设置登录用的的id
	{
		execRouter.POST("/send", exec.Send)
		//单任务
		execRouter.POST("/create-task-for-pm2", exec.CreateTaskForPM2)
		execRouter.POST("/release/:id", exec.Release)

		//多任务
		execRouter.POST("/create-task-for-pm2-v2", exec.CreateTaskForPM2V2)
		execRouter.POST("/release-v2/:id", exec.ReleaseV2)

		//scp不需要支持多任务
		execRouter.POST("/create-task-for-scp", exec.CreateTaskForSCP)
		execRouter.POST("/release-for-scp/:id", exec.ReleaseForSCP)

	}
	//权限管理
	permissionRouter := router.Group("/permission")
	permissionRouter.Use(middleware.Jwt().MiddlewareFunc())
	permissionRouter.Use(middleware.SetUserInfo()) // 自动设置登录用的的id
	{
		permissionRouter.GET("/list", permission.List)
		permissionRouter.GET("/role-permissions/:roleid", permission.RolePermissions)
		permissionRouter.POST("/create", permission.CreatePermission)
		permissionRouter.PUT("/update/:id", permission.UpdatePermission)
		permissionRouter.DELETE("/:id", permission.DeletePermission)
	}
	//角色管理
	roleRouter := router.Group("/role")
	roleRouter.Use(middleware.Jwt().MiddlewareFunc())
	roleRouter.Use(middleware.SetUserInfo()) // 自动设置登录用的的id
	{
		roleRouter.GET("/list", role.List)
		roleRouter.POST("/create", role.CreateRole)
		roleRouter.POST("/set-permission", role.SetPermission)
		roleRouter.PUT("/update/:id", role.UpdateRole)
		roleRouter.DELETE("/:id", role.DeleteRole)
	}
}
