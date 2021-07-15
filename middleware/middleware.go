package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/yuedun/zhuque/util"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/yuedun/zhuque/db"
	"github.com/yuedun/zhuque/pkg/user"

	"github.com/gin-gonic/gin"
)

//Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// 请求前
		s := string([]rune(c.HandlerName())[29:])
		log.Print("请求前 path：", c.FullPath(), s)
		c.Next()
		// 请求后
		latency := time.Since(t)
		// 获取发送的 status
		status := c.Writer.Status()
		log.Printf("请求后 path： %s 耗时： %s 状态：%d", c.FullPath(), latency, status)
	}
}

// SetUserInfo 获取用户信息中间件
// 使用c.Get("userid")获取
func SetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		if userID64, ok := claims["user_id"].(float64); ok {
			userID := int(userID64)
			log.Print("登录用户id:", userID)
			c.Keys["userid"] = userID
		}
		c.Next()
	}
}

type LoginUser struct {
	UserID      int
	UserName    string
	LoginTime   time.Time
	Permissions []string //用户拥有的权限
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
		Key:         []byte(util.Conf.JWTSecret),
		Timeout:     time.Hour * 24 * 10, //10天
		MaxRefresh:  time.Hour * 24 * 10,
		IdentityKey: identityKey,
		// 登录验证成功后存储用户信息
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*LoginUser); ok {
				return jwt.MapClaims{
					"user_id":   v.UserID,
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		//获取用户信息
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &LoginUser{
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
			userService := user.NewService(db.DB)
			userObj := user.User{
				UserName: username,
			}
			user, err := userService.GetLoginUserInfo(userObj)
			if err != nil {
				log.Println("登录查询数据", err)
				return nil, jwt.ErrFailedAuthentication
			}

			permissions := strings.Split(user.Permissions, ",")
			if user.UserName == username && user.Password == util.GetMD5(password) {
				// 返回的数据用在上面定义的PayloadFunc函数中
				log.Println(">>>><<<<", permissions)
				return &LoginUser{
					UserID:      user.ID,
					UserName:    username,
					LoginTime:   time.Now(),
					Permissions: permissions,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		// 登录以后通过token来获取用户标识，检测是否通过认证
		Authorizator: func(data interface{}, c *gin.Context) bool {
			authority := string([]rune(c.HandlerName())[29:])
			if v, ok := data.(*LoginUser); ok && v.UserName != "" {
				log.Println(">>>>>>>>>>>Authorizator", authority, v.Permissions)
				return true
			}
			log.Println("token 验证失败！")
			return false
		},
		// 获取不到token或解析token失败时如何返回信息
		Unauthorized: func(c *gin.Context, code int, message string) {
			log.Println("用户未登录，获取不到token")
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
		TokenLookup: "cookie:jwt",
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
