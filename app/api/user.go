package api

import (
	"encoding/json"
	"goj-server/app/model"
	"goj-server/app/service"
	"goj-server/global"

	"github.com/gogf/gf/net/ghttp"
)

type user struct{}

// User API管理对象
var User = new(user)

// SignUp 注册接口
func (*user) SignUp(r *ghttp.Request) {
	reqBytes := r.GetBody()
	reqData := &model.SignUpReq{}

	err := json.Unmarshal(reqBytes, reqData)
	if err != nil {
		// 返回错误
		resp, _ := json.Marshal(model.SignUpResp{
			StatusCode: global.RequestError,
			Msg:        err.Error(),
		})
		r.Response.WriteJsonExit(resp)
	}

	err, code := service.User.SignUp(reqData, r)
	if err != nil {
		// 返回错误
		resp, _ := json.Marshal(model.SignUpResp{
			StatusCode: code,
			Msg:        err.Error(),
		})
		r.Response.WriteJsonExit(resp)
	}

	// 返回提示成功数据
	resp, _ := json.Marshal(model.SignUpResp{
		StatusCode: global.RegisterSuccess,
		Msg:        global.Msg[global.RegisterSuccess],
	})
	r.Response.WriteJson(resp)
}

// LogIn 登录接口
func (*user) LogIn(r *ghttp.Request) {
	// TODO 用帐号、密码登录
	reqBytes := r.GetBody()
	reqData := &model.LogInReq{}

	json.Unmarshal(reqBytes, reqData)
	err, code := service.User.LogIn(reqData, r)

	if err != nil {
		resp, _ := json.Marshal(model.LogInResp{
			StatusCode: code,
			Msg:        err.Error(),
		})
		r.Response.WriteJsonExit(resp)
	}

	resp, _ := json.Marshal(model.LogInResp{
		StatusCode: global.RegisterSuccess,
		Msg:        global.Msg[global.LoginSuccess],
		Username:   r.Session.GetString("username"),
	})

	r.Response.WriteJson(resp)
}

// GetProfile 获得用户信息接口
func (*user) GetProfile(r *ghttp.Request) {
	// TODO 获得用户信息
	// service.GetProfile()
}

// UpdateProfile 更新用户信息接口
func (*user) UpdateProfile(r *ghttp.Request) {
	// TODO 更新用户信息
	// service.UpdateProfile()
}
