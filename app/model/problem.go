package model

type ProblemListItem struct {
	ID         int
	Name       string
	Difficulty string
}

type AddProblemReq struct {
	Name        string
	Description string
	Difficulty  string
	Testcases   []Testcase
	Tlimit      uint64
	Rlimit      uint64
}

type ProblemLimit struct {
	Tlimit int
	Rlimit int
}

type Testcase struct {
	Input  string
	Output string
}

type TestcaseDat struct {
	Input  []byte
	Output []byte
}

type ProblemContent struct {
	ID          uint64
	Name        string
	Description string
	Difficulty  string
	Tlimit      uint64
	Rlimit      uint64
}

type JudgeReq struct {
	PID      uint64
	Code     string
	Language string
}

type JudgeResult struct {
	Code   int
	Msg    string
	Tusage int
	Rusage int
}

type Submission struct {
	Language  int
	Result    int
	Tusage    int
	Rusage    int
	CreatedAt string `db:"createdAt"`
}
