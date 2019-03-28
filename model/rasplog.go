package model

import "encoding/json"

type RaspLog struct {
	*System
	StackTrace string `json:"stack_trace"`
	AppId      string `json:"app_id"`
	RaspId     string `json:"rasp_id"`
	Level      string `json:"level"`
	EventTime  string `json:"event_time"`
	Message    string `json:"message"`
	Pid        int    `json:"pid"`
	ErrorCode  int    `json:"error_code,omitempty"`
}

func (rl *RaspLog) String() string {
	b, err := json.Marshal(rl)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}
