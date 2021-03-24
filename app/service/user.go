package service

import (
	"database/sql"
	"errors"
	"goj/app/dao"
	"goj/app/model"
	"goj/global"
	"net/http"
	"strconv"
	"strings"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

// User 相关逻辑
var User = new(userService)

type userService struct{}

// SignUp 处理用户注册
func (*userService) SignUp(r *model.SignUpReq, request *ghttp.Request) (error, global.StatusCode) {
	logger := g.Log()

	// 获取用户语言
	acceptLang := strings.Split(request.GetHeader("Accept-Language"), ",")[0]
	validator := gvalid.New().I18n(acceptLang)

	// 校验请求内容
	if err := validator.CheckStruct(r, nil); err != nil {

		// 请求出错，返回400状态码
		request.Response.Status = http.StatusBadRequest
		logger.Warning(err.FirstString())
		return errors.New(err.FirstString()), global.RequestError
	}

	var username string
	// 若用户不存在则注册，否则返回错误
	if err := dao.DB.QueryRowx(
		"select username from user where username = ?", r.Username).Scan(&username); err != nil &&
		err == sql.ErrNoRows {
		// 注册新用户

		// 开始事务
		tx, err := dao.DB.Begin()
		if err != nil {
			logger.Error("创建帐号开始事务失败")

			// 服务器内部出错，返回500状态码
			request.Response.Status = http.StatusInternalServerError
			return errors.New("创建帐号时发生错误"), global.ServerError
		}

		// 新建用户
		_, err = tx.Exec(
			"insert into user(id, username, email, create_time) values (null, ?, ?, default)",
			r.Username,
			r.Email,
		)

		if err != nil {

			// 回滚事务
			tx.Rollback()
			logger.Warning("用户名或邮箱重复")

			// 请求错误，返回400状态码
			request.Response.Status = http.StatusBadRequest
			return errors.New("用户名或邮箱已被占用"), global.RequestError
		}

		// 获得新建用户的UID用于关联表的插入
		row := tx.QueryRow("select id from user where username = ?", r.Username)

		var uid int

		if err := row.Scan(&uid); err != nil {
			tx.Rollback()
			logger.Warning("找不到用户: " + r.Username)
			request.Response.Status = http.StatusInternalServerError
			return errors.New("创建帐号时发生错误"), global.ServerError
		}

		// 为新用户新建user_auth项用于登陆验证
		_, err = tx.Exec(
			"insert into user_auth(id, user_id, type, identifier, credential) values (null, ?, ?, ?, ?)",
			uid,
			"builtin",
			r.Account,
			r.Password,
		)

		if err != nil {
			tx.Rollback()
			request.Response.Status = http.StatusInternalServerError
			return errors.New("创建帐号时发生错误"), global.ServerError
		}

		// 为新用户新建个人资料
		_, err = tx.Exec(
			"insert into user_profile(id, user_id, avatar, description, update_time) values (null, ?, null, null, default)",
			uid)

		if err != nil {
			tx.Rollback()
			logger.Warning("为用户 " + strconv.Itoa(uid) + " 创建个人资料失败")
			request.Response.Status = http.StatusInternalServerError
			return errors.New("创建帐号时发生错误"), global.ServerError
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logger.Error("创建用户 " + strconv.Itoa(uid) + " 失败")

			// 提交事务失败，返回500状态码
			request.Response.Status = http.StatusInternalServerError
			return errors.New("创建帐号时发生错误"), global.ServerError
		}
	} else {
		logger.Warning("用户 " + r.Username + " 已存在")
		request.Response.Status = http.StatusBadRequest
		return errors.New("用户已存在"), global.RequestError
	}

	logger.Info("新用户注册：" + r.Username)

	return nil, global.Success
}

// LogIn 处理用户登陆
func (*userService) LogIn(r *model.LogInReq, request *ghttp.Request) (error, global.StatusCode) {

	acceptLang := request.GetQueryString("lang")
	validator := gvalid.New().I18n(acceptLang)

	if err := validator.CheckStruct(r, nil); err != nil {
		request.Response.Status = http.StatusBadRequest
		return errors.New(err.FirstString()), global.RequestError
	}

	result := dao.DB.QueryRowx(
		"select user_id from user_auth where identifier = ? AND credential = ?",
		r.Identifier,
		r.Credential)

	var uid int

	if err := result.Scan(&uid); err != nil {
		request.Response.Status = http.StatusBadRequest
		return errors.New("帐号或密码错误"), global.RequestError
	}

	var username string

	result = dao.DB.QueryRowx("select username from user where id = ?", uid)
	if err := result.Scan(&username); err != nil {
		request.Response.Status = http.StatusInternalServerError
		g.Log().Error("从数据库获取用户信息时出错")
		return errors.New("获取用户信息时出错"), global.ServerError
	}

	request.Session.Set("username", username)
	return nil, global.Success
}
