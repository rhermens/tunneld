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
	if path == "" {
		return nil, fmt.Errorf("key file path is empty")
	}

	pk, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(pk)
	if err != nil {
		return nil, err
	}

	slog.Info("Key file auth ready", "path", path)
	return ssh.PublicKeys(signer), nil
}

func authViaAgent(sockPath string) (ssh.AuthMethod, error) {
	if sockPath == "" {
		err := fmt.Errorf("SSH agent socket path is not configured and SSH_AUTH_SOCK is not set")
		return nil, err
	}

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeysCallback(agent.NewClient(conn).Signers), nil
}
