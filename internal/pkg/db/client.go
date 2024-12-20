package db

import (
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

const (
	maxIdleConnections = 10
	maxOpenConnections = 100
)

func ConfigureDB(db *sqlx.DB) {
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetMaxOpenConns(maxOpenConnections)
}

func NewDB(dsn string) *sqlx.DB {
	if dsn == "" {
		dsn = "invalid_dsn"
	}

	db, err := sqlx.Open(
		"postgres",
		dsn,
	)

	if err != nil {
		return db
	}

	ConfigureDB(db)
	fmt.Println("DB connection success!")

	return db
}
