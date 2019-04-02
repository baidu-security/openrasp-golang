package cloud

// LogReq emmm
type LogReq []map[string]interface{}

// LogResp emmm
type LogResp struct {
	Count int64 `json:"count"`
}

// Log emmm
func (c *Client) Log(t string, request []byte) error {
	body, err := c.PostRaw("/v1/agent/log/"+t, request)
	if err == nil {
		body.Close()
	}
	return err
}
