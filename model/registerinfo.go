package model

type RegisterInfo struct {
	*Server
	*Language
	Id                string `json:"id"`
	Version           string `json:"version"`
	Hostname          string `json:"hostname"`
	RegisterIp        string `json:"register_ip"`
	HeartbeatInterval uint64 `json:"heartbeat_interval"`
	RaspHome          string `json:"rasp_home"`
}

func NewRegisterInfo(server *Server, language *Language, id, hostname, raspHome, version, ip string, hbInterval uint64) *RegisterInfo {
	ri := &RegisterInfo{
		Version:           version,
		HeartbeatInterval: hbInterval,
		Id:                id,
		Hostname:          hostname,
		RaspHome:          raspHome,
		RegisterIp:        ip,
	}
	ri.Server = server
	ri.Language = language
	return ri
}
