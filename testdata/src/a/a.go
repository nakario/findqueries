package a

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

var (
	conn *sql.Conn
	db *sql.DB
	tx *sql.Tx

	dbx *sqlx.DB
	txx *sqlx.Tx
	execer sqlx.Execer = sqlx.NewDb(db, "mysql")
	execerContext sqlx.ExecerContext = sqlx.NewDb(db, "mysql")
	ext sqlx.Ext = sqlx.NewDb(db, "mysql")
	extContext sqlx.ExtContext = sqlx.NewDb(db, "mysql")
	preparer sqlx.Preparer = sqlx.NewDb(db, "mysql")
	preparerContext sqlx.PreparerContext = sqlx.NewDb(db, "mysql")
	queryer sqlx.Queryer = sqlx.NewDb(db, "mysql")
	queryerContext sqlx.QueryerContext = sqlx.NewDb(db, "mysql")
)

func testUsualUseCase() {
	rows, err := db.Query("SELECT name FROM users") // want "SELECT name FROM users"
	if err != nil {
		return
	}
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			return
		}
	}
}

func testLiteralPatterns() {
	db.Query("SELECT name FROM users") // want "SELECT name FROM users"
	db.Query(`SELECT name FROM users`) // want "SELECT name FROM users"
	db.Query( // want "SELECT\n\t\t\t\t\tname\n\t\t\t\tFROM\n\t\t\t\t\tusers\n\t"
		`	SELECT
					name
				FROM
					users
	`)
}

func testQueryPatterns() {
	db.Query(`SELECT * FROM users`) // want `SELECT \* FROM users`
	db.Query(`SELECT * FROM users WHERE id > 100`) // want `SELECT \* FROM users WHERE id > 100`
	db.Query(`ABC`) // want `ABC`
}

func testQueriers() {
	ctx := context.Background()
	var id int
	var names []string

	conn.ExecContext(ctx, `ABC`) // want `ABC`
	conn.PrepareContext(ctx, `ABC`) // want `ABC`
	conn.QueryContext(ctx, `ABC`) // want `ABC`
	conn.QueryRowContext(ctx, `ABC`) // want `ABC`
	db.Exec(`ABC`) // want `ABC`
	db.Prepare(`ABC`) // want `ABC`
	db.PrepareContext(ctx, `ABC`) // want `ABC`
	db.Query(`ABC`) // want `ABC`
	db.QueryContext(ctx, `ABC`) // want `ABC`
	db.QueryRow(`ABC`) // want `ABC`
	db.QueryRowContext(ctx, `ABC`) // want `ABC`
	tx.Exec(`ABC`) // want `ABC`
	tx.ExecContext(ctx, `ABC`) // want `ABC`
	tx.Prepare(`ABC`) // want `ABC`
	tx.PrepareContext(ctx, `ABC`) // want `ABC`
	tx.Query(`ABC`) // want `ABC`
	tx.QueryContext(ctx, `ABC`) // want `ABC`
	tx.QueryRow(`ABC`) // want `ABC`
	tx.QueryRowContext(ctx, `ABC`) // want `ABC`

	sqlx.Get(queryer, &id, `ABC`) // want `ABC`
	sqlx.GetContext(ctx, queryerContext, &id, `ABC`) // want `ABC`
	sqlx.MustExec(execer, `ABC`) // want `ABC`
	sqlx.MustExecContext(ctx, execerContext, `ABC`) // want `ABC`
	sqlx.NamedExec(ext, `ABC`, nil) // want `ABC`
	sqlx.NamedExecContext(ctx, extContext, `ABC`, nil) // want `ABC`
	sqlx.NamedQuery(ext, `ABC`, nil) // want `ABC`
	sqlx.NamedQueryContext(ctx, extContext, `ABC`, nil) // want `ABC`
	sqlx.Select(queryer, &names, `ABC`) // want `ABC`
	sqlx.SelectContext(ctx, queryerContext, &names, `ABC`) // want `ABC`
	dbx.Get(&id, `ABC`) // want `ABC`
	dbx.GetContext(ctx, &id, `ABC`) // want `ABC`
	dbx.MustExec(`ABC`) // want `ABC`
	dbx.MustExecContext(ctx, `ABC`) // want `ABC`
	dbx.NamedExec(`ABC`, nil) // want `ABC`
	dbx.NamedExecContext(ctx, `ABC`, nil) // want `ABC`
	dbx.NamedQuery(`ABC`, nil) // want `ABC`
	dbx.NamedQueryContext(ctx, `ABC`, nil) // want `ABC`
	dbx.QueryRowx(`ABC`) // want `ABC`
	dbx.QueryRowxContext(ctx, `ABC`) // want `ABC`
	dbx.Queryx(`ABC`) // want `ABC`
	dbx.QueryxContext(ctx, `ABC`) // want `ABC`
	dbx.Select(&names, `ABC`) // want `ABC`
	dbx.SelectContext(ctx, &names, `ABC`) // want `ABC`
	txx.Get(&id, `ABC`) // want `ABC`
	txx.GetContext(ctx, &id, `ABC`) // want `ABC`
	txx.MustExec(`ABC`) // want `ABC`
	txx.MustExecContext(ctx, `ABC`) // want `ABC`
	txx.NamedExec(`ABC`, nil) // want `ABC`
	txx.NamedExecContext(ctx, `ABC`, nil) // want `ABC`
	txx.NamedQuery(`ABC`, nil) // want `ABC`
	txx.PrepareNamed(`ABC`) // want `ABC`
	txx.PrepareNamedContext(ctx, `ABC`) // want `ABC`
	txx.Preparex(`ABC`) // want `ABC`
	txx.PreparexContext(ctx, `ABC`) // want `ABC`
	txx.QueryRowx(`ABC`) // want `ABC`
	txx.QueryRowxContext(ctx, `ABC`) // want `ABC`
	txx.Queryx(`ABC`) // want `ABC`
	txx.QueryxContext(ctx, `ABC`) // want `ABC`
	txx.Select(&names, `ABC`) // want `ABC`
	txx.SelectContext(ctx,  &names, `ABC`) // want `ABC`
	execer.Exec(`ABC`) // want `ABC`
	execerContext.ExecContext(ctx, `ABC`) // want `ABC`
	preparer.Prepare(`ABC`) // want `ABC`
	preparerContext.PrepareContext(ctx, `ABC`) // want `ABC`
	queryer.Query(`ABC`) // want `ABC`
	queryer.QueryRowx(`ABC`) // want `ABC`
	queryer.Queryx(`ABC`) // want `ABC`
	queryerContext.QueryContext(ctx, `ABC`) // want `ABC`
	queryerContext.QueryRowxContext(ctx, `ABC`) // want `ABC`
	queryerContext.QueryxContext(ctx, `ABC`) // want `ABC`
}
