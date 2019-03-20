package model

import "encoding/json"

type AttackLog struct {
	*AttackResult
	*Server
	*System
	*RequestInfo
	AttackParams interface{} `json:"attack_params"`
	SourceCode   []string    `json:"source_code"`
	StackTrace   string      `json:"stack_trace"`
	RaspId       string      `json:"rasp_id"`
	AppId        string      `json:"app_id"`
	ServerIp     string      `json:"server_ip"`
	EventTime    string      `json:"event_time"`
	EventType    string      `json:"event_type"`
	AttackType   string      `json:"attack_type"`
}

func (al *AttackLog) String() string {
	b, err := json.Marshal(al)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}
