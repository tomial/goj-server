package router

import (
	"goj/app/api"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	s := g.Server()
	s.Group("/user", func(group *ghttp.RouterGroup) {
		group.POST("/login", api.User.LogIn)
		group.POST("/signup", api.User.SignUp)
		group.GET("/profile", api.User.GetProfile)
		group.POST("/profile/update", api.User.UpdateProfile)
	})
}
