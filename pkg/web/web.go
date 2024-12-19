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
	http   *http.Server
	logger *zerolog.Logger
	db     *storage.Storage
}

func NewServer(db *storage.Storage, logger *zerolog.Logger) (*Server, error) {
	if db == nil {
		return nil, ErrInvalidStorage
	}
	if logger == nil {
		logger = &log.Logger
	}

	s := &Server{
		http:   &http.Server{Addr: ":8009"},
		logger: logger,
		db:     db,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger(logger))
	r.Use(middleware.WithCors)
	r.Mount("/_hotaru", s.hotaruRouter())

	s.http.Handler = r

	return s, nil
}

func (s *Server) Serve() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.http.Shutdown(context.Background())
}

func (s *Server) hotaruRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/status", s.pingHandler)
	return r
}

func (s *Server) pingHandler(w http.ResponseWriter, _ *http.Request) {
	if err := s.db.Ping(); err != nil {
		http.Error(w, "DB ERROR, PLEASE REPORT TO OP", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("OK"))
}
