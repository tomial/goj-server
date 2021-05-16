package global

type ResultCode = int

const (
	PASS                  ResultCode = iota // 0
	WRONG_ANSWER                            // 1
	RUNTIME_ERROR                           // 2
	COMPILE_ERROR                           // 3
	UNKNOWN_ERROR                           // 4
	TIME_LIMIT_EXCEEDED                     // 5
	MEMORY_LIMIT_EXCEEDED                   // 6
)

var ResultMsg = map[ResultCode]string{
	PASS:                  "提交通过",
	WRONG_ANSWER:          "答案错误",
	RUNTIME_ERROR:         "运行时出现错误",
	COMPILE_ERROR:         "编译错误",
	UNKNOWN_ERROR:         "未知错误",
	TIME_LIMIT_EXCEEDED:   "超过运行时间限制",
	MEMORY_LIMIT_EXCEEDED: "超过运行内存限制",
}
