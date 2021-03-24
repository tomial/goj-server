package model

import "goj/global"

// SignUpReq 注册请求结构体
type SignUpReq struct {
	Username string `v:"required|min-length:1"`
	Account  string `v:"required|passport"`
	Password string `v:"required|password"`
	Email    string `v:"required|email"`
}

// SignUpResp 注册响应结构体
type SignUpResp struct {
	StatusCode global.StatusCode
	Msg        string
}

// LogInReq 登录请求结构体
type LogInReq struct {
	Identifier string `v:"required|passport"`
	Credential string `v:"required|password"`
}

// LogInResp 登录响应结构体
type LogInResp struct {
	StatusCode global.StatusCode
	Msg        string
	Username   string
}
