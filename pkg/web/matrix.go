package web

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"hotaru.hana.im/server/pkg/crypto"
	"hotaru.hana.im/server/pkg/storage"
	"hotaru.hana.im/server/pkg/web/middleware"
	"hotaru.hana.im/server/pkg/web/model"
)

func (s *Server) matrixRouters() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/client/versions", clientVersionsHandler)
	r.Get("/client/v3/register/available", s.apiRegisterAvailable)
	r.With(middleware.Bind[model.RequestLogin]()).Post("/client/v3/login", s.apiLoginPost)
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

func (s *Server) apiLoginPost(w http.ResponseWriter, r *http.Request) {
	form := middleware.GetObject(r).(*model.RequestLogin)

	// bad login type
	if form.Type != model.AuthenticationTypePassword || form.Identifier.Type != model.IdentifierTypeUser {
		middleware.ErrorUnknown(w)
		return
	}

	if err := s.db.VerifyAccount(form.Identifier.User, form.Password); err != nil {
		middleware.ErrorForbidden(w)
		return
	}

	// TODO: put token into db
	ak, _ := crypto.CryptoRandomBytes(128)
	devID, _ := crypto.CryptoRandomBytes(16)

	resp := &model.ResponseLoginPost{
		UserID:      form.Identifier.User,
		AccessToken: hex.EncodeToString(ak),
		DeviceID:    hex.EncodeToString(devID),
	}
	raw, _ := json.Marshal(resp)
	w.Write(raw)
}
