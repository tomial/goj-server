package global

// 全局响应状态码
const (
	RegisterSuccess = iota
	LoginSuccess
	StatusError
)

// Msg 全局响应消息
var Msg = []string{RegisterSuccess: "Successfully registered.", LoginSuccess: "Successfully logged in", StatusError: "Error"}
