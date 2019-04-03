package orsql

import (
	openrasp "github.com/baidu-security/openrasp-golang"
	"github.com/baidu-security/openrasp-golang/common"
	"github.com/baidu-security/openrasp-golang/model"
)

type SqlErrorParam struct {
	Server  string `json:"server"`
	Query   string `json:"query"`
	ErrCode string `json:"error_code"`
	ErrMsg  string `json:"-"`
}

func NewSqlErrorParam(server, query, errCode, errMsg string) *SqlErrorParam {
	sep := &SqlErrorParam{
		Server:  server,
		Query:   query,
		ErrCode: errCode,
		ErrMsg:  errMsg,
	}
	return sep
}

func (sep *SqlErrorParam) buildPluginMessage() string {
	return sep.Server + " error " + sep.ErrCode + " detected: " + sep.ErrMsg
}

func (sep *SqlErrorParam) GetType() common.CheckType {
	return common.SqlException
}

func (sep *SqlErrorParam) GetTypeString() string {
	return common.CheckTypeToString(sep.GetType())
}

func (sep *SqlErrorParam) AttackCheck(opts ...common.AttackOption) []*model.AttackResult {
	var results []*model.AttackResult
	for _, opt := range opts {
		if opt(sep) {
			return results
		}
	}
	if sep.Server == "mysql" {
		ic := openrasp.GetAction().Get(sep.GetType())
		ar := model.NewAttackResult(model.InterceptCodeToString(ic), sep.buildPluginMessage(), "go_builtin_plugin", sep.GetTypeString(), 100)
		results = append(results, ar)
	}
	return results
}
