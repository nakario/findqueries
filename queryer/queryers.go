package main

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

var (
	conn *sql.Conn
	db *sql.DB
	tx *sql.Tx

	dbx *sqlx.DB
	txx *sqlx.Tx
	execer sqlx.Execer
	execerContext sqlx.ExecerContext
	// ext sqlx.Ext
	// extContext sqlx.ExtContext
	preparer sqlx.Preparer
	preparerContext sqlx.PreparerContext
	queryer sqlx.Queryer
	queryerContext sqlx.QueryerContext
)

var queryers = []interface{}{
	conn.ExecContext,
	conn.PrepareContext,
	conn.QueryContext,
	conn.QueryRowContext,
	db.Exec,
	db.Prepare,
	db.PrepareContext,
	db.Query,
	db.QueryContext,
	db.QueryRow,
	db.QueryRowContext,
	tx.Exec,
	tx.ExecContext,
	tx.Prepare,
	tx.PrepareContext,
	tx.Query,
	tx.QueryContext,
	tx.QueryRow,
	tx.QueryRowContext,

	// sqlx.BindNamed,
	sqlx.Get,
	sqlx.GetContext,
	// sqlx.In,
	sqlx.MustExec,
	sqlx.MustExecContext,
	// sqlx.Named,
	sqlx.NamedExec,
	sqlx.NamedExecContext,
	sqlx.NamedQuery,
	sqlx.NamedQueryContext,
	// sqlx.Preparex,
	// sqlx.PreparexContext,
	// sqlx.Rebind,
	sqlx.Select,
	sqlx.SelectContext,
	// dbx.BindNamed,
	dbx.Get,
	dbx.GetContext,
	dbx.MustExec,
	dbx.MustExecContext,
	dbx.NamedExec,
	dbx.NamedExecContext,
	dbx.NamedQuery,
	dbx.NamedQueryContext,
	dbx.QueryRowx,
	dbx.QueryRowxContext,
	dbx.Queryx,
	dbx.QueryxContext,
	// dbx.Rebind,
	dbx.Select,
	dbx.SelectContext,
	// txx.BindNamed,
	txx.Get,
	txx.GetContext,
	txx.MustExec,
	txx.MustExecContext,
	txx.NamedExec,
	txx.NamedExecContext,
	txx.NamedQuery,
	txx.PrepareNamed,
	txx.PrepareNamedContext,
	txx.Preparex,
	txx.PreparexContext,
	txx.QueryRowx,
	txx.QueryRowxContext,
	txx.Queryx,
	txx.QueryxContext,
	// txx.Rebind,
	txx.Select,
	txx.SelectContext,
	execer.Exec,
	execerContext.ExecContext,
	// ext.BindNamed,
	// ext.Exec,
	// ext.Query,
	// ext.QueryRowx,
	// ext.Queryx,
	// ext.Rebind,
	// extContext.BindNamed,
	// extContext.ExecContext,
	// extContext.QueryContext,
	// extContext.QueryRowxContext,
	// extContext.QueryxContext,
	// extContext.Rebind,
	preparer.Prepare,
	preparerContext.PrepareContext,
	queryer.Query,
	queryer.QueryRowx,
	queryer.Queryx,
	queryerContext.QueryContext,
	queryerContext.QueryRowxContext,
	queryerContext.QueryxContext,
}