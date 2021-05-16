package service

import (
	"database/sql"
	"errors"
	"fmt"
	"goj-server/app/dao"
	"goj-server/app/model"
	"goj-server/global"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

var Problems = new(problemsService)

type problemsService struct{}

func (*problemsService) GetContent(pid uint64, r *ghttp.Request) (*model.ProblemContent, error) {
	result := dao.DB.QueryRowx(`select * from problems where id = ?`, pid)
	content := &model.ProblemContent{}

	err := result.StructScan(content)
	if err != nil {
		g.Log().Warning("找不到题目内容：", pid)
		g.Log().Warning(err.Error())
		r.Response.Status = http.StatusBadRequest
		return nil, errors.New("找不到内容")
	}

	return content, nil
}

func (*problemsService) GetList(page, num int, r *ghttp.Request) ([]model.ProblemListItem, error) {
	list := make([]model.ProblemListItem, 0)
	fmt.Println(num, num*page)
	rows, err := dao.DB.Queryx(
		"select id, name, difficulty from problems LIMIT ? OFFSET ?", num, num*page)

	if err != nil {
		r.Response.Status = http.StatusBadRequest
		g.Log().Warning("获取题目列表失败")
		g.Log().Warning(err.Error())
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

func (*problemsService) Judge(jr *model.JudgeReq, r *ghttp.Request) (*model.JudgeResult, error) {
	volumeDir := g.Cfg().GetString("judger.volumeDir")
	var caseNum int
	caseRows, err := dao.DB.Queryx("select input, output from testcase where problem_id = ?", jr.PID)

	if err != nil {
		g.Log().Warningf("未找到题目#%d的测试用例", jr.PID)
		r.Response.Status = http.StatusBadRequest
		return nil, errors.New("未找到题目的测试用例")
	}

	// 尝试获得判题容器锁并写入用户代码
	switch jr.Language {
	case "go":
		fallthrough
	case "golang":
		global.GoConLock.Lock()
		os.WriteFile(volumeDir+"source.go", []byte(jr.Code), 0666)
	case "c":
		global.CConLock.Lock()
		os.WriteFile(volumeDir+"source.c", []byte(jr.Code), 0666)
	case "c++":
		fallthrough
	case "cpp":
		global.CConLock.Lock()
		os.WriteFile(volumeDir+"source.cpp", []byte(jr.Code), 0666)
	}

	// 写入输出输出
	for caseRows.Next() {
		caseNum++
		tc := &model.TestcaseDat{}
		caseRows.StructScan(tc)

		err = os.WriteFile(volumeDir+"input-"+strconv.Itoa(caseNum), tc.Input, 0644)
		if err != nil {
			g.Log().Warning(err.Error())
			g.Log().Warningf("未能获取题目#%d的所有测试用例", jr.PID)
			r.Response.Status = http.StatusInternalServerError
			return nil, errors.New("未能获取题目的所有测试用例")
		}

		err = os.WriteFile(volumeDir+"res-"+strconv.Itoa(caseNum), tc.Output, 0644)
		if err != nil {
			g.Log().Warning(err.Error())
			g.Log().Warningf("未能获取题目#%d的所有测试用例", jr.PID)
			r.Response.Status = http.StatusInternalServerError
			return nil, errors.New("未能获取题目的所有测试用例")
		}
	}

	if caseNum == 0 {
		g.Log().Warningf("未找到题目#%d的测试用例", jr.PID)
		r.Response.Status = http.StatusBadRequest
		return nil, errors.New("未找到题目的测试用例")
	}

	// 获得运行限制
	limitDat := dao.DB.QueryRowx("select tlimit, rlimit from problems where id = ?", jr.PID)
	limit := &model.ProblemLimit{}
	if err = limitDat.StructScan(limit); err != nil {
		g.Log().Warning("获取题目 %d 运行限制失败", jr.PID)
		return nil, errors.New("获取题目运行限制失败")
	}

	// 执行判题模块
	cmd := exec.Command(volumeDir+"gj", "judge", strings.ToLower(jr.Language), "-n", strconv.Itoa(caseNum), "-d", volumeDir, "-t", strconv.Itoa(limit.Tlimit), "-r", strconv.Itoa(limit.Rlimit))

	// if err = cmd.Run(); err != nil {
	// 	r.Response.Status = http.StatusInternalServerError
	// 	g.Log().Warning(err.Error())
	// 	g.Log().Warningf("判题失败，题目ID：%d", jr.PID)
	// }
	_ = cmd.Run()
	result := &model.JudgeResult{}

	// 返回结果
	if err = measure(result); err != nil {
		g.Log().Warning(err.Error())
		return nil, err
	}

	g.Log().Info("tusage:", result.Tusage)
	g.Log().Info("tlimit:", limit.Tlimit)

	switch cmd.ProcessState.ExitCode() {
	case global.PASS:
		if result.Rusage > limit.Rlimit*1024 {
			result.Code = global.MEMORY_LIMIT_EXCEEDED
			result.Msg = global.ResultMsg[global.MEMORY_LIMIT_EXCEEDED]
		} else {
			result.Code = global.PASS
			result.Msg = global.ResultMsg[global.PASS]
		}

	case global.COMPILE_ERROR:
		result.Code = global.COMPILE_ERROR
		br, _ := os.ReadFile(volumeDir + "build_result")
		result.Msg = string(br)

	case global.WRONG_ANSWER:
		result.Code = global.WRONG_ANSWER
		input, _ := os.ReadFile(volumeDir + "wa_input")
		expect, _ := os.ReadFile(volumeDir + "wa_expect")
		output, _ := os.ReadFile(volumeDir + "wa_output")
		result.Msg = fmt.Sprintf("输入：%s\n预期输出：%s\n实际输出：%s", string(input), string(expect), string(output))
		g.Log().Info(result.Msg)

	case global.RUNTIME_ERROR:
		if result.Tusage >= limit.Tlimit {
			result.Code = global.TIME_LIMIT_EXCEEDED
			result.Msg = "运行时间超过限制"
		} else {
			result.Code = global.RUNTIME_ERROR
			rte, _ := os.ReadFile(volumeDir + "rte")
			result.Msg = string(rte)
		}

	case global.UNKNOWN_ERROR:
		result.Code = global.UNKNOWN_ERROR
		result.Msg = global.ResultMsg[global.UNKNOWN_ERROR]
	}

	switch jr.Language {
	case "go":
		fallthrough
	case "golang":
		global.GoConLock.Unlock()
	case "c":
		fallthrough
	case "cpp":
		global.CConLock.Unlock()
	}

	// 将提交记录写入数据库
	uid := r.Session.GetInt("uid")
	_, err = dao.DB.Exec(`insert into submissions(id, user_id, problem_id, language, result, code, tusage, rusage, createdAt) values(null, ?, ?, ?, ?, ?, ?, ?, DEFAULT)`,
		uid,
		jr.PID,
		global.LangCode[jr.Language],
		result.Code,
		jr.Code,
		result.Tusage,
		result.Rusage)
	if err != nil {
		g.Log().Warning("保存提交记录失败")
		g.Log().Warning(err.Error())
	}

	return result, nil
}

func measure(jr *model.JudgeResult) error {
	g.Log().Info("正在获取运行数据")
	volumeDir := g.Cfg().GetString("judger.volumeDir")
	g.Log().Info("volumeDir:", volumeDir)

	path := volumeDir + "result"
	data, err := os.ReadFile(path)

	if err != nil {
		g.Log().Warning("获取运行数据失败")
		g.Log().Warning(err.Error())
		return errors.New("获取运行数据失败")
	}

	target := string(data)
	g.Log().Info("Data:", target)

	// digit := regexp.MustCompile(`(\d*[.]\d+)+|(\d+kB)+`)
	digit := regexp.MustCompile(`(\d+.\d+)+rt|(\d+)kB`)
	result := digit.FindAllString(target, -1)

	g.Log().Info("Result:", result)

	g.Log().Info(result[0][:len(result[0])-2])
	tu, err := strconv.ParseFloat(result[0][:len(result[0])-2], 32)
	if err != nil {
		g.Log().Warning("获取运行数据失败")
		return errors.New("获取运行数据失败")
	}

	jr.Tusage = int(tu * 1000)

	g.Log().Info(result[1][:len(result[1])-2])
	rt, err := strconv.ParseFloat(result[1][:len(result[1])-2], 32)
	if err != nil {
		g.Log().Warning("获取运行数据失败")
		return errors.New("获取运行数据失败")
	}

	jr.Rusage = int(rt)

	return nil
}

func (*problemsService) GetSubmissions(uid, pid int) []model.Submission {
	result, err := dao.DB.Queryx(`select language, result, tusage, rusage, createdAt from submissions
where user_id = ? and problem_id = ?`, uid, pid)
	if err != nil {
		fmt.Println(err.Error())
	}

	sub := model.Submission{}
	list := []model.Submission{}
	for result.Next() {
		err = result.StructScan(&sub)
		if err != nil {
			fmt.Println(err.Error())
		}
		list = append(list, sub)
	}

	return list
}
