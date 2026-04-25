package registry

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func (cfg *RegistryClientConfig) AddSshAuth() error {
	var methods []ssh.AuthMethod

	if cfg.SshKeyPath != "" {
		method, err := newKeyFileAuth(cfg.SshKeyPath)
		if err != nil {
			panic(err)
		}

		methods = append(methods, method)
	}

	if cfg.SshAgentSock != "" {
		method, err := newAgentAuth(cfg.SshAgentSock)
		if err != nil {
			panic(err)
		}

		methods = append(methods, method)
	}

	if len(methods) == 0 {
		return fmt.Errorf("no SSH auth methods available: provide a key file via ssh_key_path or set SSH_AUTH_SOCK")
	}

	cfg.SshConfig.Auth = methods
	return nil
}

func newKeyFileAuth(path string) (ssh.AuthMethod, error) {
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

func newAgentAuth(sockPath string) (ssh.AuthMethod, error) {
	if sockPath == "" {
		return nil, fmt.Errorf("Agent sock path is empty")
	}

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeysCallback(agent.NewClient(conn).Signers), nil
}
