package orsql

import (
	"context"
	"database/sql/driver"
	"errors"

	"github.com/baidu/openrasp"
	"github.com/baidu/openrasp/model"
)

func newConn(in driver.Conn, d *wrapDriver, dsnInfo DSNInfo) driver.Conn {
	conn := &conn{Conn: in, driver: d}
	conn.dsnInfo = dsnInfo
	conn.namedValueChecker, _ = in.(namedValueChecker)
	conn.pinger, _ = in.(driver.Pinger)
	conn.queryer, _ = in.(driver.Queryer)
	conn.queryerContext, _ = in.(driver.QueryerContext)
	conn.connPrepareContext, _ = in.(driver.ConnPrepareContext)
	conn.execer, _ = in.(driver.Execer)
	conn.execerContext, _ = in.(driver.ExecerContext)
	conn.connBeginTx, _ = in.(driver.ConnBeginTx)
	conn.connGo110.init(in)
	if in, ok := in.(driver.ConnBeginTx); ok {
		return &connBeginTx{conn, in}
	}
	return conn
}

type conn struct {
	driver.Conn
	connGo110
	driver  *wrapDriver
	dsnInfo DSNInfo

	namedValueChecker  namedValueChecker
	pinger             driver.Pinger
	queryer            driver.Queryer
	queryerContext     driver.QueryerContext
	connPrepareContext driver.ConnPrepareContext
	execer             driver.Execer
	execerContext      driver.ExecerContext
	connBeginTx        driver.ConnBeginTx
}

func (c *conn) interceptError(query string, resultError *error) {
	if *resultError == driver.ErrSkip {
		return
	}
	hit, errCode, errMsg := c.driver.errorInterceptor(resultError)
	if hit {
		sqlErrorParam := NewSqlErrorParam(c.driver.driverName, query, errCode, errMsg)
		interceptCode, _ := sqlErrorParam.AttackCheck()
		//TODO log
		if interceptCode == model.Block {
			panic(openrasp.ErrBlock)
		}
	}
}

func (c *conn) queryAttackCheck(query string) {
	sqlParam := NewSqlParam(c.driver.driverName, query)
	interceptCode, _ := sqlParam.AttackCheck()
	if interceptCode == model.Block {
		panic(openrasp.ErrBlock)
	}
}

func (c *conn) Ping(ctx context.Context) (resultError error) {
	if c.pinger == nil {
		return nil
	}
	return c.pinger.Ping(ctx)
}

func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (_ driver.Rows, resultError error) {
	if c.queryerContext == nil && c.queryer == nil {
		return nil, driver.ErrSkip
	}
	c.queryAttackCheck(query)
	defer c.interceptError(query, &resultError)

	if c.queryerContext != nil {
		return c.queryerContext.QueryContext(ctx, query, args)
	}
	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}
	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return c.queryer.Query(query, dargs)
}

func (*conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	return nil, errors.New("Query should never be called")
}

func (c *conn) PrepareContext(ctx context.Context, query string) (_ driver.Stmt, resultError error) {
	c.queryAttackCheck(query)
	defer c.interceptError(query, &resultError)
	var stmt driver.Stmt
	var err error
	if c.connPrepareContext != nil {
		stmt, err = c.connPrepareContext.PrepareContext(ctx, query)
	} else {
		stmt, err = c.Prepare(query)
		if err == nil {
			select {
			default:
			case <-ctx.Done():
				stmt.Close()
				return nil, ctx.Err()
			}
		}
	}
	if stmt != nil {
		stmt = newStmt(stmt, c, query)
	}
	return stmt, err
}

func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (_ driver.Result, resultError error) {
	if c.execerContext == nil && c.execer == nil {
		return nil, driver.ErrSkip
	}
	c.queryAttackCheck(query)
	defer c.interceptError(query, &resultError)

	if c.execerContext != nil {
		return c.execerContext.ExecContext(ctx, query, args)
	}
	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}
	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return c.execer.Exec(query, dargs)
}

func (*conn) Exec(query string, args []driver.Value) (driver.Result, error) {
	return nil, errors.New("Exec should never be called")
}

func (c *conn) CheckNamedValue(nv *driver.NamedValue) error {
	return checkNamedValue(nv, c.namedValueChecker)
}

type connBeginTx struct {
	*conn
	connBeginTx driver.ConnBeginTx
}

func (c *connBeginTx) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return c.connBeginTx.BeginTx(ctx, opts)
}
