package client

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Database interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

func New(driverName, dataSourceName string) (Database, error) {
	return sqlx.Open(driverName, dataSourceName)
}
