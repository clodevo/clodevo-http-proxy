package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/clodevo/raven-proxy/pkg/config"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(cfg config.DatabaseConfig) {
	var dsn string
	switch cfg.Type {
	case "sqlite3":
		dsn = cfg.FilePath
		ensureDir(cfg.FilePath)
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	default:
		log.Fatal("Unsupported database type")
	}

	var err error
	DB, err = sql.Open(cfg.Type, dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// It's a good practice to check if the database connection is actually alive.
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	createTables(DB)
}

func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm) // Creates the directory with permissions
	}
	return nil
}

func createTables(db *sql.DB) {
	// Create tables if they don't exist
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS admin (
        id INTEGER PRIMARY KEY,
        apikey TEXT NOT NULL UNIQUE
    );
    CREATE TABLE IF NOT EXISTS tenants (
        tenant_id CHAR(36) PRIMARY KEY,
        tenant_name TEXT NOT NULL UNIQUE
    );
    CREATE TABLE IF NOT EXISTS api_keys (
        api_key_id CHAR(36) PRIMARY KEY,
        api_key TEXT NOT NULL,
        tenant_id CHAR(36) NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
        UNIQUE(api_key, tenant_id)
    );
    `
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
}
