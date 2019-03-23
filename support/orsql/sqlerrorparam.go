package orsql

import (
	openrasp "github.com/baidu-security/openrasp-golang"
	"github.com/baidu-security/openrasp-golang/common"
	"github.com/baidu-security/openrasp-golang/gls"
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

func (sep *SqlErrorParam) AttackCheck() []*model.AttackResult {
	var results []*model.AttackResult
	ic := openrasp.GetAction().Get(common.SqlException)
	if ic != model.Ignore {
		bitMaskValue := gls.Get("whiteMask")
		bitMask, ok := bitMaskValue.(int)
		if ok && (bitMask&int(common.SqlException) == 0) {
			if sep.Server == "mysql" {
				ar := model.NewAttackResult(model.InterceptCodeToString(ic), sep.buildPluginMessage(), "go_builtin_plugin", "sql_exception", 100)
				results = append(results, ar)
			}
		}
	}
	return results
}
