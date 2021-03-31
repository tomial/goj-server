package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/frame/g"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

// NewDB 初始化数据库对象
func init() {
	link := g.Cfg().GetString("database.localhost.link")
	dbType := g.Cfg().GetString("database.type")
	logger := g.Log()

	internal, err := sqlx.Open(dbType, link)
	if err != nil {
		logger.Panic("Failed to initialize database connection: ", err)
	}
	DB = internal

	// 测试连接是否可用
	err = DB.Ping()
	if err != nil {
		logger.Panic("Failed to initialize database connection: ", err)
	}

	logger.Info("Initialized database connection.")
}
