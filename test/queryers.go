package a

import "database/sql"

var db *sql.DB

const Q = "SELECT * FROM constant"

func f1() string {
	return "SELECT * FROM f1"
}

func f2() (string, error) {
	return "SELECT * FROM f2", nil
}

func f3() (string, []interface{}) {
	return "SELECT * FROM f3 WHERE a=?, b=?", []interface{}{1, 1}
}

func queryers() {
	db.Query("SELECT * FROM interpreted_string_lit")
	db.Query(`SELECT * FROM raw_string_lit`)
	db.Query("SELECT * " + "FROM binary_operation")
	q2 := "SELECT * FROM variable"
	db.Query(q2)
	q2 = "SELECT * FROM reassigned_variable"
	db.Query(q2)
	db.Query(Q)
	q3 := "SELECT * "
	if f1() == "" {
		q3 += "FROM if"
	} else {
		q3 += "FROM else"
	}
	db.Query(q3)
	db.Query(f1())
	q4, err := f2()
	if err != nil {
		panic(err)
	}
	db.Query(q4)
	db.Query(f3())
	go db.Query("SELECT * FROM go")
	defer db.Query("SELECT * FROM defer")
}
