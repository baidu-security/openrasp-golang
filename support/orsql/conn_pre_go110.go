// +build !go1.10

package orsql

import "database/sql/driver"

type connGo110 struct{}

func (connGo110) init(in driver.Conn) {}
