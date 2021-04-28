package global

type StatusCode int

// 全局响应状态码
const (
	Success StatusCode = iota
	GenericError
	ServerError
	RequestError

	RegisterSuccess = 100
	LoginSuccess    = 200
)

// Msg 全局响应消息
var Msg = map[StatusCode]string{
	Success:         "请求成功",
	GenericError:    "发生未知错误",
	ServerError:     "服务器错误",
	RequestError:    "请求错误",
	RegisterSuccess: "注册成功",
	LoginSuccess:    "登陆成功",
}
