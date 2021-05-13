package router

import (
	"goj-server/app/api"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gsession"
)

func middleWare(r *ghttp.Request) {
	corsOptions := r.Response.DefaultCORSOptions()
	corsOptions.AllowDomain = []string{"localhost:3000"}
	corsOptions.AllowCredentials = "true"
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}

func init() {
	s := g.Server()

	s.BindMiddlewareDefault(middleWare)

	// Session 设置
	s.SetConfigWithMap(g.Map{
		"SessionMaxAge":  time.Hour * 24,
		"SessionStorage": gsession.NewStorageRedis(g.Redis()),
	})

	s.Group("/user", func(group *ghttp.RouterGroup) {
		group.POST("/login", api.User.LogIn)
		group.POST("/signup", api.User.SignUp)
		group.GET("/profile", api.User.GetProfile)
		group.POST("/profile/update", api.User.UpdateProfile)
	})

	s.Group("/problems", func(group *ghttp.RouterGroup) {
		group.GET("/", api.Problems.GetList)
		group.GET("/:id", api.Problems.Get)
		group.POST("/add", api.Problems.AddProblem)
		group.POST("/judge", api.Problems.Judge)
	})
}
