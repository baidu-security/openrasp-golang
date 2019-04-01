package cloud

import (
	"time"
)

// ReportReq emmm
type ReportReq struct {
	RaspId     string `json:"rasp_id"`
	Time       int64  `json:"time"`
	RequestSum int64  `json:"request_sum"`
}

// ReportResp emmm
type ReportResp struct{}

// Report emmm
func (c *Client) Report(requestSum int64) error {
	request := ReportReq{
		c.rasp.Id,
		int64(time.Millisecond),
		requestSum,
	}
	var response ReportResp
	return c.Post("/v1/agent/report", &request, &response)
}
