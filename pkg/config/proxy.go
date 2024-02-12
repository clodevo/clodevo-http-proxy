package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	DefaultMaxConcurrent = 512
	DefaultAddr          = ":8080"
	DefaultDNS           = ""
	DefaultTimeout       = 20 * time.Second
)

type ProxyConfig struct {
	Addr          string
	MaxConcurrent int
	DNS           []string
	Timeout       time.Duration
}

func LoadProxyConfig() *ProxyConfig {
	viper.AutomaticEnv()                                             // Read from environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")) // Replace dots and hyphens with underscores in env vars

	// Setting default values using Viper
	viper.SetDefault("proxy.addr", DefaultAddr)
	viper.SetDefault("proxy.maxConcurrent", DefaultMaxConcurrent)
	viper.SetDefault("proxy.dns", DefaultDNS)
	viper.SetDefault("proxy.timeout", DefaultTimeout)

	// Use Viper to retrieve values
	config := &ProxyConfig{
		Addr:          viper.GetString("proxy.addr"),
		MaxConcurrent: viper.GetInt("proxy.maxConcurrent"),
		DNS:           viper.GetStringSlice("proxy.dns"),
		Timeout:       viper.GetDuration("proxy.timeout"),
	}

	return config
}
