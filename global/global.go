package global

// 全局响应状态码
const (
	StatusSuccess = iota
	StatusError
)

// Msg 全局响应消息
var Msg = []string{StatusSuccess: "Success", StatusError: "Error"}
