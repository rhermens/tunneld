package registry

import (
	"fmt"
	"log/slog"
	"os"

	"golang.org/x/crypto/ssh"
)

type SshRegistryConnection struct {
	Client   *ssh.Client
	Channel  ssh.Channel
	Requests <-chan *ssh.Request
}

type RegistryClientConfig struct {
	Address      string
	SshKeyPath   string
	SshAgentSock string
	SshConfig    ssh.ClientConfig
}

func NewSshClientConfig() ssh.ClientConfig {
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	return ssh.ClientConfig{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		User:            hostName,
		Auth:            []ssh.AuthMethod{},
	}
}

func NewSshRegistryConnection(config *RegistryClientConfig, typ ConnectionType) (*SshRegistryConnection, error) {
	var err error
	connection := &SshRegistryConnection{}

	connection.Client, err = ssh.Dial("tcp", config.Address, &config.SshConfig)
	if err != nil {
		return nil, err
	}

	connection.Channel, connection.Requests, err = connection.Client.OpenChannel(string(typ), []byte{})
	if err != nil {
		return nil, err
	}

	slog.Info("Opened channel to registry", "remote", config.Address)
	return connection, nil
}

func (cfg *RegistryClientConfig) AddSshAuth() error {
	var methods []ssh.AuthMethod

	slog.Info("Configuring SSH authentication")

	method, err := authViaKeyFile(cfg.SshKeyPath)
	if err != nil {
		slog.Warn("Key file strategy unavailable", "error", err)
	} else {
		slog.Info("Key file strategy added")
		methods = append(methods, method)
	}

	method, err = authViaAgent(cfg.SshAgentSock)
	if err != nil {
		slog.Warn("SSH agent strategy unavailable", "error", err)
	} else {
		slog.Info("SSH agent strategy added")
		methods = append(methods, method)
	}

	if len(methods) == 0 {
		return fmt.Errorf("no SSH auth methods available: provide a key file via ssh_key_path or set SSH_AUTH_SOCK")
	}

	cfg.SshConfig.Auth = methods
	return nil
}

func (c *SshRegistryConnection) Close() error {
	return c.Channel.Close()
}
