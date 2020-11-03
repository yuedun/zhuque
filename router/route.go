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
	userRouter := router.Group("/user")
	// userRouter.Use(middleware.Logger())
	//user路由注册,可以给各个group加中间件
	userRouter.POST("/login", middleware.Jwt().LoginHandler)
	userRouter.GET("/refresh_token", middleware.Jwt().RefreshHandler) // 刷新token
	userRouter.GET("/logout", middleware.Jwt().LogoutHandler)
	userRouter.Use(middleware.Jwt().MiddlewareFunc())
	{
		userRouter.GET("/info/:id", middleware.Jwt().MiddlewareFunc(), user.GetUserInfo) //单独给某个路由添加中间件
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
	projectRouter.Use(middleware.Logger())
	projectRouter.Use(middleware.Jwt().MiddlewareFunc())
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
	//权限管理
	permissionRouter := router.Group("/permission")
	permissionRouter.Use(middleware.Jwt().MiddlewareFunc())
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
	{
		roleRouter.GET("/list", role.List)
		roleRouter.POST("/create", role.CreateRole)
		roleRouter.POST("/set-permission", role.SetPermission)
		roleRouter.PUT("/update/:id", role.UpdateRole)
		roleRouter.DELETE("/:id", role.DeleteRole)
	}
}
