// +build !go1.9

package orsql

import "database/sql/driver"

type stmtGo19 struct{}

func (stmtGo19) init(in driver.Stmt) {}
