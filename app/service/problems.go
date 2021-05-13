package service

import (
	"goj-server/app/dao"
	"goj-server/app/model"

	"github.com/jmoiron/sqlx"
)

var Problems = new(problemsService)

type problemsService struct{}

func (*problemsService) Get(start, end int) (*sqlx.Rows, error) {
	result, err := dao.DB.Queryx(
		"select * from problems BETWEEN ? AND ?", start, end)
	return result, err
}

func (*problemsService) AddProblem(req *model.AddProblemReq) error {
	return nil
}
