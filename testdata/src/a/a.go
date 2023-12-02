package a

import (
	"context"
	"database/sql"
	"sqlz"
)

var (
	conn *sql.Conn
	db *sql.DB
	tx *sql.Tx

	dbz *sqlz.DB
	txz *sqlz.Tx
	execer sqlz.Execer
	execerContext sqlz.ExecerContext
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
	for _, q := range [2]string{`SELECT 1 FROM array`, `SELECT 2 FROM array`} {
		db.Query(q) // want `SELECT 1 FROM array` `SELECT 2 FROM array`
	}
	array := [2]string{`SELECT 3 FROM array `, `SELECT 4 FROM array`}
	for i := range array {
		db.Query(array[i]) // want `SELECT 3 FROM array` `SELECT 4 FROM array`
	}
	for _, q := range []string{`SELECT 1 FROM slice`, `SELECT 2 FROM slice`} {
		db.Query(q) // want `SELECT 1 FROM slice` `SELECT 2 FROM slice`
	}
	slice := []string{`SELECT 3 FROM slice`, `SELECT 4 FROM slice`}
	for i := range slice {
		db.Query(slice[i]) // want `SELECT 3 FROM slice` `SELECT 4 FROM slice`
	}
	for _, q := range map[int]string {
		1: `SELECT 1 FROM map`,
		2: `SELECT 2 FROM map`,
	} {
		db.Query(q) // want `SELECT 1 FROM map` `SELECT 2 FROM map`
	}
	m := map[int]string{
		3: `SELECT 3 FROM map`,
		4: `SELECT 4 FROM map`,
	}
	for k := range m {
		db.Query(m[k]) // want `SELECT 3 FROM map` `SELECT 4 FROM map`
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
}

func testQueriers() {
	ctx := context.Background()

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

	sqlz.Exec(execer, `ABC`) // want `ABC`
	sqlz.ExecContext(ctx, execerContext, `ABC`) // want `ABC`
	dbz.Exec(`ABC`) // want `ABC`
	dbz.ExecContext(ctx, `ABC`) // want `ABC`
	txz.Exec(`ABC`) // want `ABC`
	txz.ExecContext(ctx, `ABC`) // want `ABC`
	execer.Exec(`ABC`) // want `ABC` `ABC`
	execerContext.ExecContext(ctx, `ABC`) // want `ABC` `ABC`
}
