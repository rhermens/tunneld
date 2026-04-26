package proxy

import (
	"bytes"
	"log/slog"
	"net"
	"net/http"

	"github.com/rhermens/tunneld/pkg/registry"
)

type HttpProxy struct {
	mux      *http.ServeMux
	Config   *HttpServerConfig
	Registry *registry.Registry
}

func NewHttpProxy(config *HttpServerConfig, registry *registry.Registry) HttpProxy {
	proxy := HttpProxy{
		mux:      http.NewServeMux(),
		Config:   config,
		Registry: registry,
	}

	return proxy
}

func (p *HttpProxy) Listen() error {
	for _, handlePath := range p.Config.Paths {
		p.mux.HandleFunc(handlePath, p.ForwardHandler)
		slog.Info("Registered handler", "path", handlePath)
	}

	slog.Info("Starting http server", "host", p.Config.Host, "port", p.Config.Port)
	return http.ListenAndServe(net.JoinHostPort(p.Config.Host, p.Config.Port), p.mux)
}

func (p *HttpProxy) ForwardHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Received request", "method", r.Method, "url", r.URL.String())
	var buffer bytes.Buffer
	r.Write(&buffer)

	p.Registry.FanoutBuffer(buffer.Bytes())
	w.Write([]byte("OK"))
}
