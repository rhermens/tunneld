package registry

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func (cfg *RegistryClientConfig) AddSshAuth() error {
	cfg.addKeyFile()
	cfg.addAgent()

	if len(cfg.SshConfig.Auth) == 0 {
		return fmt.Errorf("no SSH auth methods available: provide a key file via ssh_key_path or set SSH_AUTH_SOCK")
	}

	return nil
}

func (cfg *RegistryClientConfig) addAgent() {
	if cfg.SshAgentSock == "" {
		return
	}

	method, err := newAgentAuthMethod(cfg.SshAgentSock)
	if err != nil {
		slog.Warn("Failed to add ssh agent", "sock", cfg.SshAgentSock)
		return
	}

	cfg.SshConfig.Auth = append(cfg.SshConfig.Auth, method)
}

func (cfg *RegistryClientConfig) addKeyFile() {
	if cfg.SshKeyPath == "" {
		return
	}

	method, err := newKeyFileAuthMethod(cfg.SshKeyPath)
	if err != nil {
		slog.Warn("Failed to add ssh key", "path", cfg.SshKeyPath)
		return
	}

	cfg.SshConfig.Auth = append(cfg.SshConfig.Auth, method)
}

func newKeyFileAuthMethod(path string) (ssh.AuthMethod, error) {
	pk, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(pk)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}

func newAgentAuthMethod(sockPath string) (ssh.AuthMethod, error) {
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeysCallback(agent.NewClient(conn).Signers), nil
}
