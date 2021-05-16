package api

import (
	"encoding/json"
	"goj-server/app/model"
	"goj-server/app/service"
	"goj-server/global"
	"net/http"

	"github.com/gogf/gf/net/ghttp"
)

type problems struct{}

var Problems = new(problems)

func (*problems) Get(r *ghttp.Request) {
	pid := r.GetUint64("id")
	data, err := service.Problems.GetContent(pid, r)
	if err != nil {
		resp, _ := json.Marshal(model.GenericResp{
			StatusCode: global.RequestError,
			Msg:        global.Msg[global.RequestError],
		})
		r.Response.WriteJsonExit(resp)
	}

	resp, _ := json.Marshal(data)
	r.Response.WriteJson(resp)
}

func (*problems) GetList(r *ghttp.Request) {
	// 获取题目范围
	var page, num int
	page = r.GetQueryInt("page")
	num = r.GetQueryInt("num")

	result, err := service.Problems.GetList(page, num, r)

	if err != nil {
		resp, _ := json.Marshal(model.GenericResp{
			StatusCode: global.RequestError,
			Msg:        err.Error(),
		})
		r.Response.WriteJsonExit(resp)
	}

	resp, _ := json.Marshal(result)

	r.Response.WriteJson(resp)
}

func (*problems) AddProblem(r *ghttp.Request) {
	req := &model.AddProblemReq{}
	err := json.Unmarshal(r.GetBody(), req)

	if err != nil {
		resp, _ := json.Marshal(model.GenericResp{
			StatusCode: global.RequestError,
			Msg:        global.Msg[global.RequestError],
		})
		r.Response.WriteJsonExit(resp)
	}

	err = service.Problems.AddProblem(req, r)
	if err != nil {
		resp, _ := json.Marshal(model.GenericResp{
			StatusCode: global.RequestError,
			Msg:        err.Error(),
		})
		r.Response.WriteJsonExit(resp)
	}

	resp, _ := json.Marshal(model.GenericResp{
		StatusCode: global.AddProblemSucc,
		Msg:        global.Msg[global.AddProblemSucc],
	})
	r.Response.WriteJson(resp)
}

func (*problems) Judge(r *ghttp.Request) {
	body := r.GetBody()
	jr := &model.JudgeReq{}
	err := json.Unmarshal(body, &jr)

	if err != nil {
		r.Response.Status = http.StatusBadRequest
		resp, _ := json.Marshal(model.GenericResp{
			StatusCode: global.RequestError,
			Msg:        global.Msg[global.RequestError],
		})
		r.Response.WriteJsonExit(resp)
	}

	result, err := service.Problems.Judge(jr, r)
	if err != nil {
		r.Response.Status = http.StatusInternalServerError
		resp, _ := json.Marshal(model.GenericResp{
			StatusCode: global.ServerError,
			Msg:        err.Error(),
		})
		r.Response.WriteJsonExit(resp)
	}

	resp, _ := json.Marshal(result)
	r.Response.WriteJson(resp)
}

func (*problems) GetSubmissions(r *ghttp.Request) {
	uid := r.Session.GetInt("uid")
	pid := r.GetInt("pid")

	list := service.Problems.GetSubmissions(uid, pid)

	resp, _ := json.Marshal(list)
	r.Response.WriteJson(resp)
}
