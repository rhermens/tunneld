package registry

import (
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

func (c *SshRegistryConnection) Close() error {
	return c.Channel.Close()
}
