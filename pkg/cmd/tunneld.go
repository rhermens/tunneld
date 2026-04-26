package cmd

import (
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/rhermens/tunneld/pkg/proxy"
	"github.com/rhermens/tunneld/pkg/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewTunneldCmd() *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "tunneld",
		Short: "Start tunnel daemon",
		Run: func(cmd *cobra.Command, args []string) {
			NewServeCmd().Run(cmd, args)
		},
	}

	tunneldConfig()
	serveCmd.AddCommand(NewServeCmd())
	return serveCmd
}

func NewServeCmd() *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start tunneld server",
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			registry := registry.NewRegistry(registry.NewRegistryServerConfig())
			proxy := proxy.NewHttpProxy(proxy.NewHttpServerConfig(), &registry)

			wg.Go(func() {
				err, _ := registry.Listen()
				if err != nil {
					slog.Error("Failed to start registry server", "error", err)
					os.Exit(1)
				}
			})

			wg.Go(func() {
				err := proxy.Listen()
				if err != nil {
					slog.Error("Failed to start http server", "error", err)
					os.Exit(1)
				}
			})

			wg.Wait()
		},
	}

	return serveCmd
}

func tunneldConfig() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	viper.SetConfigName("tunneld")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/tunneld/")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("tunneld")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	proxy.SetConfigDefaults()
	registry.SetConfigDefaults()

	err := viper.ReadInConfig()
	if err != nil {
		slog.Warn("No config file found", "error", err)
	}
}
