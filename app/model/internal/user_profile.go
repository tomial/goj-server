package internal

import "time"

// UserProfile 用户个人信息表
type UserProfile struct {
	ID          uint32
	UID         uint32    // 用户ID
	Avatar      string    // 头像URL
	Description string    // 个人描述
	UpdateTime  time.Time // 更新时间
}
