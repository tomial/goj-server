package service

import (
	"database/sql"
	"errors"
	"fmt"
	"goj-server/app/dao"
	"goj-server/app/model"
	"net/http"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

var Problems = new(problemsService)

type problemsService struct{}

func (*problemsService) GetList(page, num int, r *ghttp.Request) ([]model.ProblemListItem, error) {
	list := make([]model.ProblemListItem, 0)
	fmt.Println(num, num*page)
	rows, err := dao.DB.Queryx(
		"select id, name, difficulty from problems LIMIT ? OFFSET ?", num, num*page)

	if err != nil {
		r.Response.Status = http.StatusBadRequest
		g.Log().Warning("获取题目列表失败")
		g.Log().Warning(err)
		return list, errors.New("获取题目列表失败")
	}

	for rows.Next() {
		item := model.ProblemListItem{}
		rows.StructScan(&item)
		fmt.Println(item.ID, item.Name, item.Difficulty)
		list = append(list, item)
	}

	return list, nil
}

func (*problemsService) AddProblem(req *model.AddProblemReq, r *ghttp.Request) error {

	cases := req.Testcases

	var pid uint64

	tx, err := dao.DB.Beginx()
	if err != nil {
		tx.Rollback()
		g.Log().Warning("开始事务时发生错误：", req.Name)
		g.Log().Warning(err.Error())
		r.Response.Status = http.StatusInternalServerError
		return errors.New("添加问题时发生错误")
	}

	result := tx.QueryRowx(`select name from problems where name = ?`, req.Name)
	if result.Scan(&pid) != sql.ErrNoRows {
		tx.Rollback()
		g.Log().Warning("重复添加问题")
		r.Response.Status = http.StatusBadRequest
		return errors.New("重复添加问题")
	}

	_, err = tx.Exec(`insert into problems(id, name, difficulty, description, tlimit, rlimit)
						values (null, ?, ?, ?, ?, ?)`,
		req.Name, req.Difficulty, req.Description, req.Tlimit, req.Rlimit)

	if err != nil {
		tx.Rollback()
		g.Log().Warning("添加问题时发生错误：", req.Name)
		g.Log().Warning(err.Error())
		r.Response.Status = http.StatusBadRequest
		return errors.New("添加问题时发生错误")
	}

	result = tx.QueryRowx(`select id from problems where name = ?`, req.Name)

	err = result.Scan(&pid)
	if err != nil {
		tx.Rollback()
		g.Log().Warning("获取问题ID时发生错误：", req.Name)
		g.Log().Warning(err.Error())
		r.Response.Status = http.StatusInternalServerError
		return errors.New("添加问题时发生错误")
	}

	for i, v := range cases {
		_, err = tx.Exec(`insert into testcase(id, problem_id, input, output) values (null, ?, ?, ?)`,
			pid, v.Input, v.Output)
		if err != nil {
			tx.Rollback()
			g.Log().Warning("添加测试用例时发生错误：", req.Name, "用例：", i)
			g.Log().Warning(err.Error())
			r.Response.Status = http.StatusBadRequest
			return errors.New("添加问题时发生错误")
		}
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		g.Log().Warning("提交添加问题事务失败：", req.Name)
		g.Log().Warning(err.Error())
		r.Response.Status = http.StatusInternalServerError
		return errors.New("添加问题时发生错误")
	}

	g.Log().Info("成功添加问题：", req.Name)

	return nil
}
