package cmd

import (
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/rhermens/tunneld/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewTunnelCmd() *cobra.Command {
	tunnelCmd := &cobra.Command{
		Use:   "tunnel",
		Short: "Listen to a tunnel",
		Run: func(cmd *cobra.Command, args []string) {
			NewListenCmd().Run(cmd, args)
		},
	}

	viper.SetConfigName("tunnel")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path.Join(os.Getenv("HOME"), ".config/tunnel"))
	viper.SetEnvPrefix("tunnel")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(".")

	tunnelCmd.Flags().String("host", "localhost", "Host to forward to")
	tunnelCmd.Flags().StringP("port", "p", "8080", "Port to forward to")
	tunnelCmd.Flags().StringP("registry", "r", ":7891", "Registry address")
	tunnelCmd.Flags().StringP("ssh-key-path", "k", path.Join(os.Getenv("HOME"), ".ssh/id_ed25519"), "Private key path for registry authentication")
	tunnelCmd.Flags().String("ssh-agent-sock", "", "SSH agent socket path (defaults to SSH_AUTH_SOCK env var)")
	viper.BindPFlag("host", tunnelCmd.Flags().Lookup("host"))
	viper.BindPFlag("port", tunnelCmd.Flags().Lookup("port"))
	viper.BindPFlag("registry", tunnelCmd.Flags().Lookup("registry"))
	viper.BindPFlag("ssh_key_path", tunnelCmd.Flags().Lookup("ssh-key-path"))
	viper.BindPFlag("ssh_agent_sock", tunnelCmd.Flags().Lookup("ssh-agent-sock"))
	viper.BindEnv("ssh_agent_sock", "SSH_AUTH_SOCK")

	err := viper.ReadInConfig()
	if err != nil {
		slog.Warn("No config file found", "error", err)
	}

	tunnelCmd.AddCommand(NewListenCmd())
	return tunnelCmd
}

func NewListenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "listen",
		Short: "Listen to a tunnel",
		Run: func(cmd *cobra.Command, args []string) {
			config := client.NewTunnelClientConfig()
			l := client.NewTunnelClient(config)
			err := l.Listen()
			if err != nil {
				slog.Error("Failed to listen to tunnel", "error", err)
				os.Exit(1)
			}

			defer l.Close()
		},
	}
}
