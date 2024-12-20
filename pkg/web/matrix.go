package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"hotaru.hana.im/server/pkg/storage"
	"hotaru.hana.im/server/pkg/web/middleware"
	"hotaru.hana.im/server/pkg/web/model"
)

func (s *Server) matrixRouters() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/client/versions", clientVersionsHandler)
	r.Get("/client/v3/register/available", s.apiRegisterAvailable)
	return r
}

func (s *Server) apiRegisterAvailable(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if err := s.db.ValidateLocalpart(username); err != nil {
		if errors.Is(err, storage.ErrUserInUse) {
			middleware.ErrorUserInUse(w)
		} else {
			middleware.ErrorUnknown(w)
		}
		return
	}

	resp := &model.ResponseRegisterAvailable{Available: true}
	raw, _ := json.Marshal(resp)
	w.Write(raw)
}
