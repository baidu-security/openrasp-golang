package orsql

import (
	"context"
	"database/sql/driver"
)

func newStmt(in driver.Stmt, conn *conn, query string) driver.Stmt {
	stmt := &stmt{
		Stmt:  in,
		conn:  conn,
		query: query,
	}
	stmt.columnConverter, _ = in.(driver.ColumnConverter)
	stmt.stmtExecContext, _ = in.(driver.StmtExecContext)
	stmt.stmtQueryContext, _ = in.(driver.StmtQueryContext)
	stmt.namedValueChecker, _ = in.(namedValueChecker)
	if stmt.namedValueChecker == nil {
		stmt.namedValueChecker = conn.namedValueChecker
	}
	return stmt
}

type stmt struct {
	driver.Stmt
	conn  *conn
	query string

	columnConverter   driver.ColumnConverter
	namedValueChecker namedValueChecker
	stmtExecContext   driver.StmtExecContext
	stmtQueryContext  driver.StmtQueryContext
}

func (s *stmt) queryAttackCheck() {
	s.conn.queryAttackCheck(s.query)
}

func (s *stmt) interceptError(resultError *error) {
	s.conn.interceptError(s.query, resultError)
}

func (s *stmt) ColumnConverter(idx int) driver.ValueConverter {
	if s.columnConverter != nil {
		return s.columnConverter.ColumnConverter(idx)
	}
	return driver.DefaultParameterConverter
}

func (s *stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (_ driver.Result, resultError error) {
	s.queryAttackCheck()
	defer s.interceptError(&resultError)
	if s.stmtExecContext != nil {
		return s.stmtExecContext.ExecContext(ctx, args)
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
	return s.Exec(dargs)
}

func (s *stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (_ driver.Rows, resultError error) {
	s.queryAttackCheck()
	defer s.interceptError(&resultError)
	if s.stmtQueryContext != nil {
		return s.stmtQueryContext.QueryContext(ctx, args)
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
	return s.Query(dargs)
}

func (s *stmt) CheckNamedValue(nv *driver.NamedValue) error {
	return checkNamedValue(nv, s.namedValueChecker)
}
