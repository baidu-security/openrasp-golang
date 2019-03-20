package ormysql

import (
	"strconv"

	"github.com/go-sql-driver/mysql"
)

func MysqlInterceptError(err *error) (bool, string, string) {
	if driverErr, ok := (*err).(*mysql.MySQLError); ok {
		errCode := strconv.Itoa(int(driverErr.Number))
		switch errCode {
		case "1045", "1060", "1062", "1064", "1105", "1367", "1690":
			return true, errCode, driverErr.Message
		}
	}
	return false, "", ""
}
