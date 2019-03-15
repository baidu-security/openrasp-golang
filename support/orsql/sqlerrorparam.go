package orsql

import "github.com/baidu/openrasp/model"

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

func (sep *SqlErrorParam) AttackCheck() (model.InterceptCode, *model.AttackResult) {
	if sep.Server == "mysql" {
		switch sep.ErrCode {
		case "1060", "1062", "1064", "1105", "1367", "1690":
			ar := model.NewAttackResult("block", sep.buildPluginMessage(), "go_builtin_plugin", "sql_exception", 100)
			return model.Block, ar
		}
	}
	return model.Ignore, nil
}
