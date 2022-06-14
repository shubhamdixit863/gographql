package db

import (
	"database/sql"
	"log"
)

var Db *sql.DB

// for a given DSN.
func OpenDB(dsn string) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	Db = db
}
