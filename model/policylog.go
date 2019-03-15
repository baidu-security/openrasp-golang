package model

import "encoding/json"

type PolicyLog struct {
	*PolicyResult
	*Server
	*System
	PolicyParams interface{} `json:"policy_params"`
	SourceCode   []string    `json:"source_code"`
	StackTrace   string      `json:"stack_trace"`
	RaspId       string      `json:"rasp_id"`
	AppId        string      `json:"app_id"`
	EventTime    string      `json:"event_time"`
}

func (pl *PolicyLog) String() string {
	b, err := json.Marshal(pl)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}
