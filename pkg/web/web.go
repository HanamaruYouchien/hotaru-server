package web

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"hotaru.hana.im/server/pkg/web/middleware"
)

type Server struct {
	HTTP   *http.Server
	Logger *zerolog.Logger
}

func NewServer(logger *zerolog.Logger) *Server {
	if logger == nil {
		logger = &log.Logger
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger(logger))
	r.Mount("/_hotaru", hotaruRouter())

	return &Server{
		HTTP:   &http.Server{Addr: ":8009", Handler: r},
		Logger: logger,
	}
}

func (s *Server) Serve() error {
	return s.HTTP.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.HTTP.Shutdown(context.Background())
}

func hotaruRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	return r
}
