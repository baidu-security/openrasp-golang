package model

type PolicyChecker interface {
	PolicyCheck() (InterceptCode, *PolicyResult)
}

type AttackChecker interface {
	AttackCheck() (InterceptCode, *AttackResult)
}
