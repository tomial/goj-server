package internal

// UserAuth 用户登录验证表
type UserAuth struct {
	ID         uint32
	UID        uint32
	AuthType   string
	Identifier string
	credential string
}
