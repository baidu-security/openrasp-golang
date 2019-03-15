package orsql

import (
	"database/sql/driver"
	"reflect"
	"strings"
)

func ExtractName(d driver.Driver) string {
	t := reflect.TypeOf(d)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Name() {
	case "SQLiteDriver":
		return "sqlite3"
	case "MySQLDriver":
		return "mysql"
	case "Driver":
		if strings.HasSuffix(t.PkgPath(), "github.com/lib/pq") {
			return "postgresql"
		}
	}
	return "unsupported"
}
