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

const Q = "SELECT a FROM constant"

func f1() string {
	return `SELECT a FROM f1`
}

func f2() (string, error) {
	return `SELECT a FROM f2`, nil
}

func f3() (string, []interface{}) {
	return "SELECT a FROM f3 WHERE b=?, c=?", []interface{}{1, 1}
}

type myString string
type stringAlias = string

func testComplexCalls() {
	db.Query((`SELECT a FROM paren`)) // want `SELECT a FROM paren`
	db.Query(`SELECT a ` + `FROM binary_operation`) // want `SELECT a FROM binary_operation`
	db.Query(string("SELECT a FROM type_conversion1")) // want `SELECT a FROM type_conversion1`
	db.Query(string(myString("SELECT a FROM type_conversion2"))) // want `SELECT a FROM type_conversion2`
	db.Query(stringAlias(myString("SELECT a FROM type_conversion3"))) // want `SELECT a FROM type_conversion3`
	q2 := `SELECT a FROM variable`
	db.Query(q2) // want `SELECT a FROM variable`
	q2 = `SELECT a FROM reassigned_variable`
	db.Query(q2) // want `SELECT a FROM reassigned_variable`
	db.Query(Q) // want `SELECT a FROM constant`
	q3 := "SELECT a "
	if f1() == "" {
		q3 += "FROM if"
	} else {
		q3 += "FROM else"
	}
	db.Query(q3) // want `SELECT a FROM if` `SELECT a FROM else`
	for _, q := range []string{`SELECT 1 FROM slice`, `SELECT 2 FROM slice`} {
		db.Query(q) // want `SELECT 1 FROM slice` `SELECT 2 FROM slice`
	}
	db.Query(f1()) // want `SELECT a FROM f1`
	q4, err := f2()
	if err != nil {
		panic(err)
	}
	db.Query(q4) // want `SELECT a FROM f2`
	db.Query(f3()) // want `SELECT a FROM f3`
	go db.Query(`SELECT a FROM go`) // want `SELECT a FROM go`
	defer db.Query("SELECT a FROM defer") // want `SELECT a FROM defer`

	inQuery, inArgs, err := sqlx.In("SELECT a FROM sqlx_in WHERE 1 in () AND a=?", []int{1, 2, 3}, 4)
	if err != nil {
		panic(err)
	}
	a := struct{a, b int}{}
	dbx.Get(&a, inQuery, inArgs...) // want `SELECT a FROM sqlx_in WHERE 1 in \(\) AND a=?`
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
