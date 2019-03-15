package common

import (
	"net"

	"github.com/baidu/openrasp/model"
	"github.com/baidu/openrasp/utils"
)

type Globals struct {
	Hostname string
	RaspId   string
	HttpAddr string

	Language      *model.Language
	System        *model.System
	Server        *model.Server
	ContextServer *model.ContextServer
}

func NewGlobals(rootPath string) *Globals {
	hostname := utils.GetHostname()
	nic, _ := getServerNic()
	raspId := calculateRaspId(hostname, rootPath)
	g := Globals{
		Hostname:      hostname,
		RaspId:        raspId,
		Language:      model.NewLanguage(),
		System:        model.NewSystem(nic, hostname),
		ContextServer: model.NewContextServer(),
	}
	return &g
}

func (g *Globals) SetServer(server *model.Server) {
	g.Server = server
}

func (g *Globals) SetHttpAddr(httpAddr string) {
	g.HttpAddr = httpAddr
}

func getServerNic() ([]model.Nic, error) {
	var nics []model.Nic
	ifaces, err := net.Interfaces()
	if err != nil {
		return nics, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPNet:
				if !v.IP.IsLoopback() {
					ipString := v.IP.String()
					nic := model.Nic{
						Name: i.Name,
						Ip:   ipString,
					}
					nics = append(nics, nic)
				}
			}
		}
	}
	return nics, nil
}

func calculateRaspId(hostname, path string) string {
	var raspString string
	macAddrs := utils.GetMacAddrs()
	for _, v := range macAddrs {
		raspString += v
	}
	raspString += hostname
	raspString += path
	return utils.GetMd5Hash(raspString)
}
