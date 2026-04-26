package keystore

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/google/go-github/v75/github"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

type GitHubConfig struct {
	Organization string
	Token        string
}

type GitHubKeystore struct {
	authorizedKeys map[string]bool
	config         GitHubConfig
	client         *github.Client
	ticker         *time.Ticker
	mu             sync.Mutex
	wg             sync.WaitGroup
}

func NewGitHubKeystore() *GitHubKeystore {
	config := GitHubConfig{
		Organization: viper.GetString("registry.ssh.github.organization"),
		Token:        viper.GetString("registry.ssh.github.token"),
	}

	keystore := &GitHubKeystore{
		config: config,
		client: github.NewClient(nil).WithAuthToken(config.Token),
		ticker: time.NewTicker(viper.GetDuration("registry.ssh.github.refresh_interval")),
	}

	keystore.wg.Go(func() {
		keystore.Refresh()
		for range keystore.ticker.C {
			keystore.Refresh()
		}
	})

	return keystore
}

func (k *GitHubKeystore) Refresh() {
	slog.Info("Fetching authorized keys from GitHub organization", "organization", k.config.Organization)
	var authorizedKeys []string

	members, _, err := k.client.Organizations.ListMembers(context.Background(), k.config.Organization, nil)
	if err != nil {
		slog.Error("Failed to list organization members", "error", err)
		return
	}

	for _, member := range members {
		keys, _, err := k.client.Users.ListKeys(context.Background(), member.GetLogin(), nil)
		if err != nil {
			slog.Error("Failed to list user keys", "user", member.GetLogin(), "error", err)
			continue
		}

		for _, key := range keys {
			authorizedKeys = append(authorizedKeys, key.GetKey())
		}
	}

	k.mu.Lock()
	k.authorizedKeys = AuthorizedKeysFromStrings(authorizedKeys)
	k.mu.Unlock()
}

func (k *GitHubKeystore) ContainsKey(key ssh.PublicKey) bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	return k.authorizedKeys[string(key.Marshal())]
}
