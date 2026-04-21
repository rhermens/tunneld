package registry

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func authViaKeyFile(path string) (ssh.AuthMethod, error) {
	slog.Info("Attempting key file auth", "path", path)

	if path == "" {
		return nil, fmt.Errorf("key file path is empty")
	}

	pk, err := os.ReadFile(path)
	if err != nil {
		slog.Warn("Key file auth failed: could not read file", "path", path, "error", err)
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(pk)
	if err != nil {
		slog.Warn("Key file auth failed: could not parse private key", "path", path, "error", err)
		return nil, err
	}

	slog.Info("Key file auth ready", "path", path)
	return ssh.PublicKeys(signer), nil
}

func authViaAgent(sockPath string) (ssh.AuthMethod, error) {
	slog.Info("Attempting SSH agent auth", "path", sockPath)

	if sockPath == "" {
		err := fmt.Errorf("SSH agent socket path is not configured and SSH_AUTH_SOCK is not set")
		slog.Warn("SSH agent auth failed", "error", err)
		return nil, err
	}

	slog.Info("Connecting to SSH agent socket", "path", sockPath)
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		slog.Warn("SSH agent auth failed: could not connect to socket", "path", sockPath, "error", err)
		return nil, err
	}

	slog.Info("SSH agent auth ready", "path", sockPath)
	return ssh.PublicKeysCallback(agent.NewClient(conn).Signers), nil
}
