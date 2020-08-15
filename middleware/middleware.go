package middleware

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/yuedun/zhuque/db"
	"github.com/yuedun/zhuque/pkg/user"

	"github.com/gin-gonic/gin"
)

//Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// 设置 example 变量
		c.Set("example", "12345")

		// 设置cookie
		c.SetCookie("user_name", "yuedun", 3600, "/", "localhost", true, true)

		// 请求前

		c.Next()

		// 请求后
		latency := time.Since(t)
		log.Print("耗时：", latency)

		// 获取发送的 status
		status := c.Writer.Status()
		log.Println("状态：", status)
	}
}

// 权限校验
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.Abort() //不继续执行
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
			})
			return
		}
		log.Println(">>>", token)
		if token == "" {
			log.Print("权限验证未通过")
			c.Abort() //不继续执行
			c.JSON(http.StatusForbidden, gin.H{
				"message": "权限验证未通过",
			})
			return
		} else {
			c.Next() //如果通过中间件需要调用Next，使其继续执行下一个func
		}
	}
}

type User struct {
	UserName  string
	LoginTime time.Time
}
type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Jwt() *jwt.GinJWTMiddleware {
	var identityKey = "username"
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour * 24 * 3,
		MaxRefresh:  time.Hour * 24 * 3,
		IdentityKey: identityKey,
		// 登录验证成功后存储用户信息
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		//获取用户信息
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims[identityKey].(string),
			}
		},
		// 首次通过用户名密码登录认证
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password
			userService := user.NewService(db.SQLLite)
			userObj := user.User{
				UserName: username,
			}
			user, err := userService.GetUserInfo(userObj)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			if (user.UserName == "admin" && user.Password == password) || (username == "test" && password == "test") {
				// 返回的数据用在上面定义的PayloadFunc函数中
				return &User{
					UserName:  username,
					LoginTime: time.Now(),
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		// 登录以后通过token来获取用户标识，检测是否通过认证
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && v.UserName == "test" {
				return true
			}

			return false
		},
		// 获取不到token或解析token失败时如何返回信息
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// 获取jwt token的方法，从header中获取，从query中获取，从cookie中获取
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
		//LoginHandler,LogoutHandler等handler中间件会默认提供，但其返回的数据格式并不一定符合项目规范，也可以在此处自定义，像上面Unauthorized这样
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	return authMiddleware
}
