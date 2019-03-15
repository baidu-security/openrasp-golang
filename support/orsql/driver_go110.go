// +build go1.10

package orsql

import (
	"context"
	"database/sql/driver"
)

func (d *wrapDriver) OpenConnector(name string) (driver.Connector, error) {
	if dc, ok := d.Driver.(driver.DriverContext); ok {
		oc, err := dc.OpenConnector(name)
		if err != nil {
			return nil, err
		}
		return &wrapConnector{oc.Connect, d, name}, nil
	}
	connect := func(context.Context) (driver.Conn, error) {
		return d.Driver.Open(name)
	}
	return &wrapConnector{connect, d, name}, nil
}

type wrapConnector struct {
	connect func(context.Context) (driver.Conn, error)
	driver  *wrapDriver
	name    string
}

func (d *wrapConnector) Connect(ctx context.Context) (driver.Conn, error) {
	dsnInfo := d.driver.dsnParser(d.name)
	conn, err := d.connect(ctx)
	if err != nil {
		return nil, err
	}
	return newConn(conn, d.driver, dsnInfo), nil
}

func (d *wrapConnector) Driver() driver.Driver {
	return d.driver
}
