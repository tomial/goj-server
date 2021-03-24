package shared

import "time"

type ContextUser struct {
	ID         uint32
	Username   string
	Email      string
	CreateTime time.Time
	// Status     uint8 TODO 邮箱验证
}
