package ormysql

import (
	"github.com/baidu/openrasp/support/orsql"
	"github.com/go-sql-driver/mysql"
)

func init() {
	orsql.Register("mysql", &mysql.MySQLDriver{}, orsql.DSNParserWrap(MysqlParseDSN), orsql.ErrorInterceptorWrap(MysqlInterceptError))
}
