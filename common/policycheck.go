package common

import (
	"github.com/baidu-security/openrasp-golang/model"
)

type PolicyChecker interface {
	PolicyCheck() (model.InterceptCode, *model.PolicyResult)
}
