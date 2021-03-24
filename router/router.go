package router

import (
	"goj/app/api"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gsession"
)

func init() {
	s := g.Server()

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
}
