package model

type AddProblemReq struct {
	Name        string
	Description string
	Difficulty  string
	Testcases   []testcase
	Tlimit      uint64
	Rlimit      uint64
}

type testcase struct {
	Input  string
	Output string
}