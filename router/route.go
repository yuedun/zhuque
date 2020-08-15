package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yuedun/zhuque/middleware"
	"github.com/yuedun/zhuque/pkg/exec"
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
		userRouter.GET("/index", user.Index)
		//userRouter.POST("/login", user.Login)
		userRouter.POST("/login", middleware.Jwt().LoginHandler)
		userRouter.GET("/refresh_token", middleware.Jwt().RefreshHandler) // 刷新token
		userRouter.GET("/logout", middleware.Jwt().LogoutHandler)
		userRouter.GET("/info/:id", middleware.Auth(), user.GetUserInfo) //单独给某个路由添加中间件
		userRouter.GET("/users-by-sql/:id", user.GetUserInfoBySql)
		userRouter.POST("/", user.CreateUser)
		userRouter.PUT("/update/:id", user.UpdateUser)
		userRouter.DELETE("/del/:id", user.DeleteUser)
	}
	//发送命令路由注册
	execRouter := router.Group("/exec")
	{
		execRouter.GET("/", exec.Index)
		execRouter.POST("/send", exec.Exec)
	}
}
