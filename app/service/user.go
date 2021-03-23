package service

import (
	"database/sql"
	"errors"
	"goj/app/dao"
	"goj/app/model"
	"strconv"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gvalid"
)

// User 相关逻辑
var User = new(userService)

type userService struct{}

// SignUp 处理用户注册
func (*userService) SignUp(r *model.SignUpReq) error {
	// 校验请求内容
	if err := gvalid.CheckStruct(r, nil); err != nil {
		return errors.New(err.FirstString())
	}

	logger := g.Log()
	var username string
	// 若用户不存在则注册，否则返回错误
	if err := dao.DB.QueryRowx(
		"select username from user where username = ?", r.Username).Scan(&username); err != nil &&
		err == sql.ErrNoRows {
		// 注册新用户

		tx, err := dao.DB.Begin()
		if err != nil {
			logger.Info("创建帐号开始事物失败")
			return errors.New("创建帐号时发生错误")
		}

		_, err = tx.Exec(
			"insert into user(id, username, email, create_time) values (null, ?, ?, default)",
			r.Username,
			r.Email,
		)

		if err != nil {
			tx.Rollback()
			logger.Info("用户名或邮箱重复")
			return errors.New("创建帐号时发生错误")
		}

		row := tx.QueryRow("select id from user where username = ?", r.Username)

		var uid int

		if err := row.Scan(&uid); err != nil {
			tx.Rollback()
			logger.Info("找不到用户: " + r.Username)
			return errors.New("创建帐号时发生错误")
		}

		_, err = tx.Exec(
			"insert into user_auth(id, user_id, type, identifier, credential) values (null, ?, ?, ?, ?)",
			uid,
			"builtin",
			r.Account,
			r.Password,
		)

		if err != nil {
			tx.Rollback()
			return errors.New("创建帐号时发生错误")
		}

		_, err = tx.Exec(
			"insert into user_profile(id, user_id, avatar, description, update_time) values (null, ?, null, null, default)",
			uid)

		if err != nil {
			tx.Rollback()
			logger.Info("为用户 " + strconv.Itoa(uid) + " 创建个人资料失败")
			return errors.New("创建帐号时发生错误")
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logger.Info("创建用户 " + strconv.Itoa(uid) + " 失败")
			return errors.New("创建帐号时发生错误")
		}
	} else {
		logger.Info("用户 " + r.Username + " 已存在")
		return errors.New("用户已存在")
	}

	logger.Info("新用户注册：" + r.Username)

	return nil
}

// LogIn 处理用户登陆
func (*userService) LogIn(r *model.LogInReq) error {
	if err := gvalid.CheckStruct(r, nil); err != nil {
		return errors.New(err.FirstString())
	}

	return nil
}
