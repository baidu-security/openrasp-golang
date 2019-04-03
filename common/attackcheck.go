package common

import (
	"github.com/baidu-security/openrasp-golang/model"
)

type AttackOption func(AttackChecker) bool

type AttackChecker interface {
	AttackCheck(opts ...AttackOption) []*model.AttackResult
	GetType() CheckType
	GetTypeString() string
}
