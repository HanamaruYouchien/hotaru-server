package web

import "github.com/go-chi/chi/v5"

func (s *Server) matrixRouters() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/client/versions", clientVersionsHandler)
	return r
}
