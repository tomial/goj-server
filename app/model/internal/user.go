package internal

import "time"

// User 用户表
type User struct {
	ID         uint32
	Username   string
	Email      string
	CreateTime time.Time
	// Status     uint8 TODO 邮箱验证
}
