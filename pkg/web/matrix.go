package web

import (
	"encoding/hex"
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
	r.Get("/client/v3/login", s.apiLoginGet)
	r.With(middleware.Bind[model.RequestLogin]()).Post("/client/v3/login", s.apiLoginPost)
	r.With(middleware.Bind[model.RequestRegister]()).Post("/client/v3/register", s.apiRegister)
	return r
}

func (s *Server) apiRegisterAvailable(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if err := s.db.ValidateLocalpart(username); err != nil {
		switch {
		case errors.Is(err, storage.ErrUserInUse):
			middleware.ErrorUserInUse(w)
		case errors.Is(err, storage.ErrInvalidUsername):
			middleware.ErrorInvalidUsername(w)
		default:
			middleware.ErrorUnknown(w)
		}
		return
	}

	resp := &model.ResponseRegisterAvailable{Available: true}
	middleware.RenderJSON(w, resp)
}

func (s *Server) apiLoginPost(w http.ResponseWriter, r *http.Request) {
	form := middleware.GetObject(r).(*model.RequestLogin)

	// bad login type
	if form.Type != model.AuthenticationTypePassword || form.Identifier.Type != model.IdentifierTypeUser {
		middleware.ErrorUnknown(w)
		return
	}

	if err := s.db.VerifyAccount(form.Identifier.User, form.Password); err != nil {
		if errors.Is(err, storage.ErrAccountNotExist) || errors.Is(err, storage.ErrPasswordNotCorrect) {
			middleware.ErrorForbiddenMsg(w, "username or password not correct")
		} else {
			middleware.ErrorForbidden(w)
		}
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
	middleware.RenderJSON(w, resp)
}

func (s *Server) apiLoginGet(w http.ResponseWriter, _ *http.Request) {
	resp := &model.ResponseLoginGet{
		Flows: []model.LoginFlow{
			{Type: model.AuthenticationTypePassword},
		},
	}
	middleware.RenderJSON(w, resp)
}

func (s *Server) apiRegister(w http.ResponseWriter, r *http.Request) {
	form := middleware.GetObject(r).(*model.RequestRegister)
	if err := s.db.CreateAccount(form.Username, form.Password, storage.AccountTypeUser); err != nil {
		switch {
		case errors.Is(err, storage.ErrPasswordTooWeak):
			middleware.ErrorWeakPassword(w)
		case errors.Is(err, storage.ErrUserInUse):
			middleware.ErrorUserInUse(w)
		case errors.Is(err, storage.ErrInvalidUsername):
			middleware.ErrorInvalidUsername(w)
		default:
			middleware.ErrorUnknown(w)
		}
		return
	}

	resp := &model.ResponseRegister{
		UserID: "@" + form.Username + ":" + s.domain,
	}
	middleware.RenderJSON(w, resp)
}
