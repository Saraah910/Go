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

	connStr := "host=localhost port=5432 user=sakshi.aherkar dbname=kuberpa sslmode=disable"
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
	CREATE TABLE IF NOT EXISTS clusters (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		provisioner TEXT NOT NULL CHECK (provisioner IN ('aws', 'azure', 'gcp', 'nutanix', 'vmware')),
		region TEXT NOT NULL,
		workspace TEXT NOT NULL DEFAULT 'default',
		kubeconfig TEXT NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		user_id INTEGER,
		status TEXT NOT NULL DEFAULT 'Pending',
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	if _, err := DB.Exec(kubeClusterTable); err != nil {
		return fmt.Errorf("error creating kubeclusters table: %w", err)
	}

	InfraTable := `
	CREATE TABLE IF NOT EXISTS infra (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		provider TEXT NOT NULL CHECK (provider IN ('aws', 'gcp', 'azure', 'nutanix')),
		config JSONB NOT NULL,
		is_default BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now(),
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	if _, err := DB.Exec(InfraTable); err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}

	return nil
}
