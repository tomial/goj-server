package shared

import (
	"github.com/gogf/gf/net/ghttp"
)

type Context struct {
	Session *ghttp.Session
	User    *ContextUser
}
