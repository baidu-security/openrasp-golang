package orsql

import (
	"encoding/json"

	openrasp "github.com/baidu-security/openrasp-golang"
	"github.com/baidu-security/openrasp-golang/common"
	"github.com/baidu-security/openrasp-golang/model"
	v8 "github.com/baidu-security/openrasp-v8/go"
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

func (sp *SqlParam) Bytes() []byte {
	b, _ := json.Marshal(sp)
	return b
}

func (sp *SqlParam) GetType() common.CheckType {
	return common.Sql
}

func (sp *SqlParam) GetTypeString() string {
	return common.CheckTypeToString(sp.GetType())
}

func (sp *SqlParam) AttackCheck(opts ...common.AttackOption) []*model.AttackResult {
	var ars []*model.AttackResult
	for _, opt := range opts {
		if opt(sp) {
			return ars
		}
	}
	if openrasp.RequestInfoAvailable() {
		resultBytes := v8.Check(sp.GetTypeString(), sp.Bytes(), openrasp.DefaultContextGetters(), openrasp.GetGeneral().GetInt("plugin.timeout.millis"))
		var ms []map[string]interface{}
		err := json.Unmarshal(resultBytes, &ms)
		if err == nil {
			for _, m := range ms {
				ars = append(ars, model.NewAttackResultFromMap(m))
			}
		}
	}
	return ars
}
