package cloud

// Rasp emmm
type Rasp struct {
	Id                string            `json:"id"`
	AppId             string            `json:"app_id"`
	Version           string            `json:"version"`
	HostName          string            `json:"hostname"`
	RegisterIp        string            `json:"register_ip"`
	Language          string            `json:"language"`
	LanguageVersion   string            `json:"language_version"`
	ServerType        string            `json:"server_type"`
	ServerVersion     string            `json:"server_version"`
	RaspHome          string            `json:"rasp_home"`
	PluginVersion     string            `json:"plugin_version"`
	HeartbeatInterval int64             `json:"heartbeat_interval"`
	Online            bool              `json:"online"`
	LastHeartbeatTime int64             `json:"last_heartbeat_time"`
	RegisterTime      int64             `json:"register_time"`
	Environ           map[string]string `json:"environ"`
}

// Plugin emmm
type Plugin struct {
	Id              string                 `json:"id"`
	AppId           string                 `json:"app_id"`
	Name            string                 `json:"name"`
	UploadTime      int64                  `json:"upload_time"`
	Version         string                 `json:"version"`
	Md5             string                 `json:"md5"`
	Content         string                 `json:"plugin,omitempty"`
	AlgorithmConfig map[string]interface{} `json:"algorithm_config"`
}
