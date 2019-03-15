package ormysql

import (
	"strconv"

	"github.com/go-sql-driver/mysql"
)

func MysqlInterceptError(err *error) (bool, string, string) {
	if driverErr, ok := (*err).(*mysql.MySQLError); ok {
		return true, strconv.Itoa(int(driverErr.Number)), driverErr.Message
	}
	return false, "", ""
}
