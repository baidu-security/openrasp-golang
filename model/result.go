package model

type InterceptCode int

const (
	Block InterceptCode = iota
	Log
	Ignore
)

type AttackResult struct {
	PluginMessage    string `json:"plugin_message"`
	PluginConfidence uint64 `json:"plugin_confidence"`
	PluginAlgorithm  string `json:"plugin_algorithm"`
	PluginName       string `json:"plugin_name"`
	InterceptState   string `json:"intercept_state"`
}

func NewAttackResult(state, message, algorithm, name string, confidence uint64) *AttackResult {
	ar := &AttackResult{
		PluginMessage:    message,
		PluginConfidence: confidence,
		PluginAlgorithm:  algorithm,
		PluginName:       name,
		InterceptState:   state,
	}
	return ar
}

type PolicyResult struct {
	EventType string `json:"event_type"`
	Message   string `json:"message"`
	PolicyId  uint64 `json:"policy_id"`
}

func NewPolicyResult(message string, policyId uint64) *PolicyResult {
	pr := &PolicyResult{
		EventType: "security_policy",
		Message:   message,
		PolicyId:  policyId,
	}
	return pr
}
