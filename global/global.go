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
	Success:         "Success",
	GenericError:    "Error",
	ServerError:     "Server Error",
	RequestError:    "Request Error,",
	RegisterSuccess: "Successfully registered.",
	LoginSuccess:    "Successfully logged in",
}
