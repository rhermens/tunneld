package proxy

import (
	"github.com/spf13/viper"
)

type HttpServerConfig struct {
	Host  string
	Port  string
	Paths []string
}

func NewHttpServerConfig() *HttpServerConfig {
	cfg := &HttpServerConfig{
		Host:  viper.GetString("http.host"),
		Port:  viper.GetString("http.port"),
		Paths: make([]string, 0),
	}
	viper.UnmarshalKey("http.paths", &cfg.Paths)

	return cfg
}

func SetConfigDefaults() {
	viper.SetDefault("http.host", "0.0.0.0")
	viper.SetDefault("http.port", "8000")
	viper.SetDefault("http.paths", []string{"/{path...}"})
}
