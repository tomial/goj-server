package api

import (
	"encoding/json"
	"goj-server/app/model"
	"goj-server/app/service"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type problems struct{}

var Problems = new(problems)

func (*problems) GetList(r *ghttp.Request) {
	// 获取题目范围
	var start, end int
	start = r.GetQueryInt("start")
	end = r.GetQueryInt("end")

	if start == 0 || end == 0 {
		start = 1
		end = 10
	}

	result, err := service.Problems.Get(start, end)
	if err != nil {
		g.Log().Warning(err.Error())
	}

	r.Response.Write(result)
}

func (*problems) Get(r *ghttp.Request) {

}

func (*problems) AddProblem(r *ghttp.Request) {
	req := &model.AddProblemReq{}
	err := json.Unmarshal(r.GetBody(), req)
	if err != nil {
		g.Log().Warning(err.Error())
	}

	service.Problems.AddProblem(req)

	g.Log().Info(req)
	r.Response.Writeln("OK")
}

func (*problems) Judge(r *ghttp.Request) {
}
