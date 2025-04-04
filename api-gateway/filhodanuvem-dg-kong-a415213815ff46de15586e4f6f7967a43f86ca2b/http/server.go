package http

import (
	"log/slog"
	"net/http"

	"github.com/devgymbr/kong"
	"github.com/devgymbr/kong/plugin"
)

type Server struct {
	Config *kong.Config
	Client ForwardClient
}

func NewServer(config *kong.Config) *Server {
	return &Server{
		Config: config,
		Client: ForwardClient{Client: http.DefaultClient, Config: config},
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service, route := s.Config.FindServiceRoute(r)
	if service == nil || route == nil {
		slog.Debug("no service found", slog.String("method", r.Method), slog.String("url", r.URL.Path))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		err := s.Client.ForwardRequest(service.URL, w, r)

		if err != nil {
			slog.Error("could not forward request", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	for j := range service.Plugins {
		p := service.Plugins[j]
		middleware, err := plugin.FindMiddleware(p.Name)
		if err != nil {
			slog.Error("unable to find middleware for plugin",
				slog.String("plugin", p.Name),
				slog.String("error", err.Error()),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		f = middleware(p, f)
	}

	f(w, r)
}
