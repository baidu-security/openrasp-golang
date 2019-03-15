package ormysql

import (
	"net"

	"github.com/baidu/openrasp/support/orsql"
	"github.com/go-sql-driver/mysql"
)

func ensureHavePort(addr string) string {
	if _, _, err := net.SplitHostPort(addr); err != nil {
		return net.JoinHostPort(addr, "3306")
	}
	return addr
}

func MysqlParseDSN(name string) orsql.DSNInfo {
	dsnInfo := orsql.DSNInfo{
		ConnectionString: name,
	}
	cfg, err := mysql.ParseDSN(name)
	if err != nil {
		return dsnInfo
	}
	dsnInfo.Database = cfg.DBName
	dsnInfo.User = cfg.User

	if cfg.Net == "" {
		cfg.Net = "tcp"
	}

	if cfg.Addr == "" {
		switch cfg.Net {
		case "tcp":
			cfg.Addr = "127.0.0.1:3306"
		case "unix":
			cfg.Addr = "/tmp/mysql.sock"
		default:
			//skip
		}

	} else if cfg.Net == "tcp" {
		cfg.Addr = ensureHavePort(cfg.Addr)
	}

	switch cfg.Net {
	case "tcp":
		host, port, err := net.SplitHostPort(cfg.Addr)
		if err == nil {
			dsnInfo.Hostname = host
			dsnInfo.Port = port
		}
	case "unix":
		dsnInfo.Socket = cfg.Addr
	default:
		//skip
	}
	return dsnInfo
}
