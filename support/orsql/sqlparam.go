package orsql

import "github.com/baidu/openrasp/model"

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
	//TODO call js
	return model.Ignore, nil
}
