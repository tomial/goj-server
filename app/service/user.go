package service

import (
	"database/sql"
	"errors"
	"goj/app/dao"
	"goj/app/model"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gvalid"
)

// User 相关逻辑
var User = new(userService)

type userService struct {
}

func (*userService) SignUp(r *model.SignUpReq) error {
	// 校验请求内容
	if err := gvalid.CheckStruct(r, nil); err != nil {
		return errors.New(err.FirstString())
	}

	logger := g.Log()
	var username string
	// 若用户不存在则注册，否则返回错误
	if err := dao.DB.QueryRowx("select username from user where username = ?", r.Username).Scan(&username); err != nil &&
		err == sql.ErrNoRows {
		// 注册新用户
		logger.Info("Registered user.")
	} else {
		return errors.New("用户已存在")
	}

	return nil
}
