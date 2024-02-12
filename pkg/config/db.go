package config

import (
	"strings"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Type     string
	FilePath string // For SQLite
	Host     string // For PostgreSQL and MySQL
	Port     int    // For PostgreSQL and MySQL
	User     string // For PostgreSQL and MySQL
	Password string // For PostgreSQL and MySQL
	DBName   string // For PostgreSQL and MySQL
}

func LoadDatabaseConfig() DatabaseConfig {
	// Setting default values
	viper.SetDefault("database.Type", "sqlite3")
	viper.SetDefault("database.FilePath", "/opt/clodevo/data.db")
	viper.SetDefault("database.Port", 3306) // Default for MySQL, adjust if using PostgreSQL
	viper.SetDefault("database.DBName", "clod-proxy")

	viper.AutomaticEnv()                                             // Read from environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")) // Replace dots and hyphens with underscores in env vars

	return DatabaseConfig{
		Type:     viper.GetString("database.Type"),
		FilePath: viper.GetString("database.FilePath"),
		Host:     viper.GetString("database.Host"),
		Port:     viper.GetInt("database.Port"),
		User:     viper.GetString("database.User"),
		Password: viper.GetString("database.Password"),
		DBName:   viper.GetString("database.DBName"),
	}
}
