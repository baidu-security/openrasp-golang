package cloud

// LogReq emmm
type LogReq []map[string]interface{}

// LogResp emmm
type LogResp struct {
	Count int64 `json:"count"`
}

// Log emmm
func (c *Client) Log(t string, request []map[string]interface{}) (int64, error) {
	var response LogResp
	err := c.Post("/v1/agent/log/"+t, &request, &response)
	return response.Count, err
}
