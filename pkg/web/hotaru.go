package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) hotaruRouters() *chi.Mux {
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
