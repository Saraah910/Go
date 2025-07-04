package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
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

	workspaceTable := `
	CREATE TABLE IF NOT EXISTS workspaces (
		id TEXT NOT NULL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		owner_id INTEGER NOT NULL,
		created_at TIMESTAMPTZ DEFAULT now(),
		members JSONB DEFAULT '[]',
		roles JSONB DEFAULT '[]',
		cluster_count INTEGER DEFAULT 0,
		cloud_providers JSONB DEFAULT '[]',
		apps_count INTEGER DEFAULT 0,
		monitoring_enabled BOOLEAN DEFAULT FALSE,
		logging_enabled BOOLEAN DEFAULT FALSE,
		tags JSONB DEFAULT '{}',
		FOREIGN KEY(owner_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	if _, err := DB.Exec(workspaceTable); err != nil {
		return fmt.Errorf("error creating workspaces table: %w", err)
	}

	kubeClusterTable := `
	CREATE TABLE IF NOT EXISTS clusters (
		id TEXT NOT NULL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		provisioner TEXT NOT NULL CHECK (provisioner IN ('aws', 'azure', 'gcp', 'nutanix', 'vmware')),
		region TEXT NOT NULL,
		kubeconfig TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now(),
		user_id INTEGER,
		workspace_id TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'Pending',
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL,
		FOREIGN KEY(workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
	);`

	if _, err := DB.Exec(kubeClusterTable); err != nil {
		return fmt.Errorf("error creating clusters table: %w", err)
	}

	infraTable := `
	CREATE TABLE IF NOT EXISTS infra (
		id TEXT NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		provider TEXT NOT NULL CHECK (provider IN ('aws', 'gcp', 'azure', 'nutanix')),
		config JSONB NOT NULL,
		is_default BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now(),
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL
	);`

	if _, err := DB.Exec(infraTable); err != nil {
		return fmt.Errorf("error creating infra table: %w", err)
	}

	return nil
}

func GetUUID() (string, error) {
	return uuid.New().String(), nil
}
