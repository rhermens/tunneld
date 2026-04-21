package registry

import (
	"fmt"
	"log/slog"
	"os"

	"golang.org/x/crypto/ssh"
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
