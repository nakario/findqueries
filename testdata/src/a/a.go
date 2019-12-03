package a

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
