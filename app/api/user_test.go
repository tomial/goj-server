package api

import (
	"encoding/json"
	"goj/app/model"
	"goj/global"
	"testing"

	"github.com/gogf/gf/frame/g"
)

func TestSignUp(t *testing.T) {
	client := g.Client()

	cases := []struct {
		name      string
		req       model.SignUpReq
		wantError bool
	}{
		{
			name: "正常请求",
			req: model.SignUpReq{
				Username: "Mudai",
				Account:  "Mudaidai",
				Password: "Mudaidai",
				Email:    "Mudai@Mudai.com",
			},
			wantError: false,
		},
		{
			name:      "空请求",
			req:       model.SignUpReq{},
			wantError: true,
		},
		{
			name: "不完整请求",
			req: model.SignUpReq{
				Username: "1",
				Password: "123456",
			},
			wantError: true,
		},
	}

	for _, c := range cases {

		req, _ := json.Marshal(c.req)
		t.Log(string(req))
		resp, err := client.Post("http://localhost:8080/user/signup", req)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		signUpResp := model.SignUpResp{}
		err = json.Unmarshal(resp.ReadAll(), &signUpResp)

		if err != nil {
			t.Log(err)
			t.FailNow()
		} else if signUpResp.StatusCode == global.StatusError && !c.wantError {
			t.Log(c.name, "FAILED")
			t.Log(signUpResp.Msg)
			t.FailNow()
		} else if signUpResp.StatusCode == global.RegisterSuccess && c.wantError {
			t.Log(c.name, "FAILED")
			t.Log(signUpResp.Msg)
			t.FailNow()
		}

		t.Log(c.name, "PASSED")
	}
}
