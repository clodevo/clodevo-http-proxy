package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type GitSyncConfig struct {
	RepoURL      string
	BranchName   string
	Username     string
	Password     string
	RepoPath     string
	SyncInterval time.Duration
}

func LoadGitSyncConfig() *GitSyncConfig {
	viper.SetDefault("git-acl.repo-path", "/opt/clodevo/acl") // Default repository path
	viper.SetDefault("git-acl.sync-interval", "1m")           // Default sync interval, e.g., "5m" for 5 minutes

	viper.AutomaticEnv()                                             // Read from environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")) // Replace dots and hyphens with underscores in env vars

	return &GitSyncConfig{
		RepoURL:      viper.GetString("git-acl.repo-url"),
		BranchName:   viper.GetString("git-acl.branch-name"),
		Username:     viper.GetString("git-acl.username"),
		Password:     viper.GetString("git-acl.password"),
		RepoPath:     viper.GetString("git-acl.repo-path"),
		SyncInterval: viper.GetDuration("git-acl.sync-interval"),
	}
}
