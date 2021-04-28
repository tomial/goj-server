package api

import "github.com/gogf/gf/net/ghttp"

type problems struct{}

var Problems = new(problems)

func (*problems) GetAll(r *ghttp.Request) {

}
