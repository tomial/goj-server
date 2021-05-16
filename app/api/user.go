package api

import (
	"encoding/json"
	"fmt"
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
			Msg:        global.Msg[global.RequestError],
			Uid:        -1,
		})
		r.Response.WriteJsonExit(resp)
	}

	err, code := service.User.SignUp(reqData, r)
	if err != nil {
		// 返回错误
		resp, _ := json.Marshal(model.SignUpResp{
			StatusCode: code,
			Msg:        err.Error(),
			Uid:        -1,
		})
		r.Response.WriteJsonExit(resp)
	}

	// 返回提示成功数据
	resp, _ := json.Marshal(model.SignUpResp{
		StatusCode: global.RegisterSuccess,
		Msg:        global.Msg[global.RegisterSuccess],
		Uid:        r.Session.GetInt("uid"),
	})
	r.Response.WriteJson(resp)
}

// LogIn 登录接口
func (*user) LogIn(r *ghttp.Request) {
	reqBytes := r.GetBody()
	reqData := &model.LogInReq{}

	err := json.Unmarshal(reqBytes, reqData)
	if err != nil {
		resp, _ := json.Marshal(model.LogInResp{
			StatusCode: global.RequestError,
			Msg:        global.Msg[global.RequestError],
			Username:   "",
		})
		r.Response.WriteJsonExit(resp)
	}

	fmt.Println(reqData)

	err, code := service.User.LogIn(reqData, r)

	if err != nil {
		resp, _ := json.Marshal(model.LogInResp{
			StatusCode: code,
			Msg:        err.Error(),
			Username:   "",
		})
		r.Response.WriteJsonExit(resp)
	}

	resp, _ := json.Marshal(model.LogInResp{
		StatusCode: global.LoginSuccess,
		Msg:        global.Msg[global.LoginSuccess],
		Username:   r.Session.GetString("username"),
	})

	r.Response.WriteJson(resp)
}

func (*user) GetRole(r *ghttp.Request) {
	uid := r.Session.GetInt("uid")

	role := service.User.GetRole(uid)

	r.Response.Write(role)
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
