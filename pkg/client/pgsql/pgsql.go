package pgsql

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const driverName = "pgx"

func NewPostgres(url string) (*sql.DB, error) {
	conn, err := sql.Open(driverName, url)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return conn, nil
}
