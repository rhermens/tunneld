package registry

import (
	"fmt"
	"log/slog"

	"github.com/rhermens/tunneld/pkg/registry/keystore"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

type SshConfig struct {
	AllowInternalUnauthenticated bool
	Keystores                    keystore.Keystores
	HostKeyPath                  string
	SshConfig                    *ssh.ServerConfig
}

type RegistryServerConfig struct {
	Host string
	Port string
	Ssh  *SshConfig
}

func SetConfigDefaults() {
	viper.SetDefault("registry.host", "0.0.0.0")
	viper.SetDefault("registry.port", "7891")
	viper.SetDefault("registry.ssh.host_key_path", ".ssh/id_ed25519")
	viper.SetDefault("registry.ssh.authorized_keys", []string{})
	viper.SetDefault("registry.ssh.github.organization", nil)
	viper.SetDefault("registry.ssh.github.token", nil)
}

func NewRegistryServerConfig() *RegistryServerConfig {
	return &RegistryServerConfig{
		Host: viper.GetString("registry.host"),
		Port: viper.GetString("registry.port"),
		Ssh:  newSshConfig(),
	}
}

func newSshConfig() *SshConfig {
	keystores := keystore.LoadKeystores()

	config := &SshConfig{
		AllowInternalUnauthenticated: viper.GetBool("registry.ssh.allow_internal_unauthenticated"),
		Keystores:                    keystores,
		HostKeyPath:                  viper.GetString("registry.ssh.host_key_path"),
		SshConfig: &ssh.ServerConfig{
			NoClientAuth:      false,
			PublicKeyCallback: newPublicKeyCallback(keystores),
		},
	}

	config.SshConfig.AddHostKey(EnsureHostKey(config.HostKeyPath))
	return config
}

func newPublicKeyCallback(ks keystore.Keystores) func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	return func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
		if ks.ContainsKey(key) {
			return &ssh.Permissions{
				Extensions: map[string]string{
					"pubkey-fp": ssh.FingerprintSHA256(key),
				},
			}, nil
		}

		slog.Warn("Unauthorized public key", "user", conn.User(), "remote", conn.RemoteAddr(), "key", key.Type())
		return nil, fmt.Errorf("unauthorized public key for %q", conn.User())
	}
}
