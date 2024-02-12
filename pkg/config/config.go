package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	DatabaseConfig DatabaseConfig
	GitSyncConfig  GitSyncConfig
	ProxyConfig    ProxyConfig
	AdminAPIKey    string
	ACLDataPath    string
	AdminAddr      string
	LogLevel       string
}

func LoadAppConfig() *AppConfig {

	viper.SetConfigName("config") // Adjust based on your actual config file
	viper.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // Or wherever your config file is located

	// Set defaults
	viper.SetDefault("acl-data-path", "/opt/clodevo/acl/tenants")
	viper.SetDefault("admin-api-key", "")
	viper.SetDefault("admin-addr", ":9090") // Default admin server address
	viper.SetDefault("log-Level", "info")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")) // Replace dots and hyphens with underscores in env vars

	// Attempt to read the configuration
	if err := viper.ReadInConfig(); err != nil {
		// Instead of failing fatally, log the error and proceed with defaults
		log.Printf("Warning: Error reading config file, using default configuration: %s", err)
	}

	return &AppConfig{
		DatabaseConfig: LoadDatabaseConfig(),
		ProxyConfig:    *LoadProxyConfig(),   // Load proxy config
		GitSyncConfig:  *LoadGitSyncConfig(), // Load Git sync config
		AdminAPIKey:    viper.GetString("admin-api-key"),
		ACLDataPath:    viper.GetString("acl-data-path"),
		AdminAddr:      viper.GetString("admin-addr"),
		LogLevel:       viper.GetString("log-Level"),
	}
}
