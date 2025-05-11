package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")
	if err != nil {
		panic("Cannot connect to database.")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createUserTable()
}

func createUserTable() {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)
	`
	_, err := DB.Exec(usersTable)
	if err != nil {
		panic("Cannot create table.")
	}

	kubeClusterTable := `
	CREATE TABLE IF NOT EXISTS kubeclusters(
		cluster_id INTEGER PRIMARY KEY AUTOINCREMENT,
		cluster_name TEXT NOT NULL,
		provisioner TEXT NOT NULL,
		kubeconfig_path TEXT NOT NULL,
		status TEXT NOT NULL, 
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`
	_, err = DB.Exec(kubeClusterTable)
	if err != nil {
		panic("Cannot create table.")
	}
}
