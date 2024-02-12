package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/clodevo/raven-proxy/pkg/config" // Adjust the import path to where your config package is located
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	http "github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Sync initializes the synchronization process for a Git repository based on the configuration.
func Sync(cfg *config.GitSyncConfig) {
	ticker := time.NewTicker(cfg.SyncInterval)
	go func() {
		for range ticker.C {
			err := syncRepo(cfg)
			if err != nil {
				fmt.Printf("Error syncing Git repository: %s\n", err)
			} else {
				fmt.Println("Git repository synced successfully")
			}
		}
	}()
}

// syncRepo handles the actual synchronization logic for a Git repository.
func syncRepo(g *config.GitSyncConfig) error {
	_, err := os.Stat(g.RepoPath)
	if os.IsNotExist(err) {
		// Clone the repository if it does not exist
		_, err := git.PlainClone(g.RepoPath, false, &git.CloneOptions{
			URL: g.RepoURL,
			Auth: &http.BasicAuth{
				Username: g.Username, // These credentials should be secured
				Password: g.Password,
			},
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", g.BranchName)),
			SingleBranch:  true,
		})
		if err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	} else {
		// Open the existing repository
		r, err := git.PlainOpen(g.RepoPath)
		if err != nil {
			return fmt.Errorf("failed to open repository: %w", err)
		}

		w, err := r.Worktree()
		if err != nil {
			return fmt.Errorf("failed to get worktree: %w", err)
		}

		// Pull the latest changes
		err = w.Pull(&git.PullOptions{
			RemoteName: "origin",
			Auth: &http.BasicAuth{
				Username: g.Username,
				Password: g.Password,
			},
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", g.BranchName)),
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return fmt.Errorf("failed to pull repository: %w", err)
		}
	}

	return nil
}
