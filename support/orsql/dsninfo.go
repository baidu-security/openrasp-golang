package orsql

type DSNInfo struct {
	Database         string `json:"-"`
	Hostname         string `json:"hostname"`
	User             string `json:"username"`
	Socket           string `json:"socket"`
	Port             string `json:"port"`
	ConnectionString string `json:"connectionString"`
}

func (dsnInfo *DSNInfo) IsHighPrivileged(driverName string) bool {
	item := driverName + ":" + dsnInfo.User
	switch item {
	case
		"mysql:root",
		"pgsql:postgres":
		return true
	}
	return false
}
