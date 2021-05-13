package model

import "goj-server/global"

type GenericResp struct {
	StatusCode global.StatusCode
	Msg        string
}
