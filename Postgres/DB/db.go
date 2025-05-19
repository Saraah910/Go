package DB

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func IntiDB() {
	var err error
	connStr := "host=localhost port=5432 user=postgres password=secret dbname=api sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic("cannot connect to DB")
	}
	DB.SetMaxIdleConns(10)
	DB.SetMaxOpenConns(5)

	createTable()
}

func createTable() {
	createUsers := `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)
	`
	_, err := DB.Exec(createUsers)
	if err != nil {
		fmt.Print(err)
		panic("Could not create table")
	}
}
