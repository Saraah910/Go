package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error

	connStr := "host=localhost port=5432 user=postgres password=secret dbname=api sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	if err = createTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}

func createTables() error {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL,
		org_name TEXT NOT NULL,
		org_department TEXT NOT NULL,
		city_location TEXT NOT NULL,
		permission TEXT NOT NULL
	);`

	if _, err := DB.Exec(userTable); err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}

	kubeClusterTable := `
	CREATE TABLE IF NOT EXISTS kubeclusters (
		cluster_id SERIAL PRIMARY KEY,
		cluster_name TEXT NOT NULL,
		provisioner TEXT NOT NULL,
		kubeconfig_path TEXT NOT NULL,
		status TEXT,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	if _, err := DB.Exec(kubeClusterTable); err != nil {
		return fmt.Errorf("error creating kubeclusters table: %w", err)
	}

	return nil
}
