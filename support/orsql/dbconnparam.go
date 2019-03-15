package orsql

import (
	"github.com/baidu/openrasp"
	"github.com/baidu/openrasp/model"
)

type DbConnectionParam struct {
	*DSNInfo
	Server string `json:"server"`
}

func NewDbConnectionParam(dsnInfo *DSNInfo, server string) *DbConnectionParam {
	dcp := &DbConnectionParam{
		Server: server,
	}
	dcp.DSNInfo = dsnInfo
	return dcp
}

func (dcp *DbConnectionParam) PolicyCheck() (model.InterceptCode, *model.PolicyResult) {
	enforcePolicy := openrasp.GetGeneral().GetBool("security.enforce_policy")
	if (dcp.DSNInfo).IsHighPrivileged(dcp.Server) {
		msg := "Database security - Connecting to a " + dcp.Server + " instance using the high privileged account: " + dcp.DSNInfo.User
		if len(dcp.DSNInfo.Socket) != 0 {
			msg += " (via unix domain socket)"
		}
		pr := model.NewPolicyResult(msg, 3006)
		is := model.Log
		if enforcePolicy {
			is = model.Block
		}
		return is, pr
	}
	return model.Ignore, nil
}
