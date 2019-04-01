package cloud

// RegisterReq emmm
type RegisterReq struct {
	*Rasp
}

// RegisterResp emmm
type RegisterResp struct {
	*Rasp
}

// Register emmm
func (c *Client) Register(id, rasphome, hostname, language, languageversion string, heartbeatinterval int64) error {
	request := RegisterReq{
		&Rasp{
			Id:                id,
			RaspHome:          rasphome,
			HostName:          hostname,
			Language:          language,
			LanguageVersion:   languageversion,
			HeartbeatInterval: heartbeatinterval,
			Version:           "1.0.0",
		},
	}
	var response RegisterResp
	if err := c.Post("/v1/agent/rasp", &request, &response); err != nil {
		return err
	}
	if response.Rasp != nil {
		c.rasp = *response.Rasp
	}
	return nil
}
