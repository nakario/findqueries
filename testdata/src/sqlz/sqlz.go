package sqlz

import "context"

type DB struct{}

type Tx struct{}

type Execer interface {
	Exec(query string) error
}

type ExecerContext interface {
	ExecContext(ctx context.Context, query string) error
}

func NewDb(db *DB, driver string) *DB {
	return nil
}

func Exec(e Execer, query string) error {
	return e.Exec(query)
}

func ExecContext(ctx context.Context, e ExecerContext, query string) error {
	return e.ExecContext(ctx, query)
}

func (db *DB) Exec(query string) error {
	return nil
}

func (db *DB) ExecContext(ctx context.Context, query string) error {
	return nil
}

func (tx *Tx) Exec(query string) error {
	return nil
}

func (tx *Tx) ExecContext(ctx context.Context, query string) error {
	return nil
}
