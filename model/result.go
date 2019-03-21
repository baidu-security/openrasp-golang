package model

type InterceptCode int

const (
	Block InterceptCode = iota
	Log
	Ignore
)

func InterceptCodeToString(code InterceptCode) string {
	switch code {
	case Block:
		return "block"
	case Log:
		return "log"
	default:
		return "ignore"
	}
}

func InterceptStringToCode(key string) InterceptCode {
	switch key {
	case "block":
		return Block
	case "ignore":
		return Ignore
	default:
		return Log
	}
}

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

func NewAttackResultFromMap(m map[string]interface{}) *AttackResult {
	ar := &AttackResult{}
	for k, v := range m {
		switch k {
		case "action":
			if state, ok := v.(string); ok {
				ar.InterceptState = state
			}
		case "message":
			if message, ok := v.(string); ok {
				ar.PluginMessage = message
			}
		case "confidence":
			if confidence, ok := v.(int); ok {
				ar.PluginConfidence = uint64(confidence)
			}
		case "name":
			if name, ok := v.(string); ok {
				ar.PluginName = name
			}
		case "algorithm":
			if algorithm, ok := v.(string); ok {
				ar.PluginAlgorithm = algorithm
			}
		default:
		}
	}
	return ar
}

func (ar *AttackResult) GetInterceptState() InterceptCode {
	return InterceptStringToCode(ar.InterceptState)
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
