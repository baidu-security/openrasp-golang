package model

type Nic struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

type System struct {
	Hostname string `json:"server_hostname"`
	Nic      []Nic  `json:"server_nic"`
}

func NewSystem(nics []Nic, hostname string) *System {
	system := &System{
		Hostname: hostname,
		Nic:      nics,
	}
	return system
}
