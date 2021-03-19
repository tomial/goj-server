package api

import (
	"encoding/json"
	"goj/app/model"
	"goj/app/service"
	"goj/global"

	"github.com/gogf/gf/net/ghttp"
)

type user struct{}

// User API管理对象
var User = new(user)

// SignUp 注册接口
func (*user) SignUp(r *ghttp.Request) {
	// TODO 用用户名、帐号、密码、邮箱注册
	reqBytes := r.GetBody()
	reqData := model.SignUpReq{}

	err := json.Unmarshal(reqBytes, &reqData)
	if err != nil {
		// 返回错误
		resp, _ := json.Marshal(model.SignUpResp{
			StatusCode: global.StatusError,
			Msg:        err.Error(),
		})
		r.Response.WriteJsonExit(resp)
	}

	err = service.User.SignUp(&reqData)
	if err != nil {
		// 返回错误
		resp, _ := json.Marshal(model.SignUpResp{
			StatusCode: global.StatusError,
			Msg:        err.Error(),
		})
		r.Response.WriteJsonExit(resp)
	}

	// 返回提示成功数据
	resp, _ := json.Marshal(model.SignUpResp{
		StatusCode: global.StatusSuccess,
		Msg:        global.Msg[global.StatusSuccess],
	})
	r.Response.WriteJson(resp)
}

// LogIn 登录接口
func (*user) LogIn(r *ghttp.Request) {
	// TODO 用帐号、密码登录
	// service.LogIn(identifier, credential)
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
