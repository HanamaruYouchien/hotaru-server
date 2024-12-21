package web

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"hotaru.hana.im/server/pkg/storage"
	"hotaru.hana.im/server/pkg/web/middleware"
)

var ErrInvalidStorage = errors.New("invalid storage")

type Server struct {
	domain string
	http   *http.Server
	logger *zerolog.Logger
	db     *storage.Storage
}

const maxBodySize = 10 * 1024 * 1024

func NewServer(domain string, db *storage.Storage, logger *zerolog.Logger) (*Server, error) {
	if db == nil {
		return nil, ErrInvalidStorage
	}
	if logger == nil {
		logger = &log.Logger
	}

	s := &Server{
		domain: domain,
		http:   &http.Server{Addr: ":8009"},
		logger: logger,
		db:     db,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger(logger), middleware.WithHeaders, middleware.MaxBodyLength(maxBodySize))
	r.Mount("/_matrix", s.matrixRouters())
	r.Mount("/_hotaru", s.hotaruRouters())

	s.http.Handler = r

	return s, nil
}

func (s *Server) Serve() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.http.Shutdown(context.Background())
}
