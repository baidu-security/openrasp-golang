package model

type Server struct {
	ServerType    string `json:"server_type"`
	ServerVersion string `json:"server_version"`
}

func NewServer(serverType, serverVersion string) *Server {
	server := &Server{
		ServerType:    serverType,
		ServerVersion: serverVersion,
	}
	return server
}
