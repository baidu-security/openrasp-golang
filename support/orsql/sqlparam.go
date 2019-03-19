package orsql

import (
	"github.com/baidu-security/openrasp-golang/common"
	"github.com/baidu-security/openrasp-golang/gls"
	"github.com/baidu-security/openrasp-golang/model"
)

type SqlParam struct {
	Query  string `json:"query"`
	Server string `json:"server"`
}

func NewSqlParam(server, query string) *SqlParam {
	sp := &SqlParam{
		Server: server,
		Query:  query,
	}
	return sp
}

func (sp *SqlParam) AttackCheck() (model.InterceptCode, *model.AttackResult) {
	bitMaskValue := gls.Get("whiteMask")
	bitMask, ok := bitMaskValue.(int)
	if ok && (bitMask&int(common.Sql) == 0) {
		//TODO call js
	}
	return model.Ignore, nil
}
