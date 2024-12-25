package web

import (
	"github.com/go-chi/chi/v5"
	"hotaru.hana.im/server/pkg/web/middleware"
	"hotaru.hana.im/server/pkg/web/model"
)

func (s *Server) matrixRouters() *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/client", s.matrixClientRouters())
	return r
}

func (s *Server) matrixClientRouters() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/versions", clientVersionsHandler)
	r.Get("/v3/register/available", s.apiRegisterAvailable)
	r.Get("/v3/login", s.apiLoginGet)
	r.With(middleware.Bind[model.RequestLogin]()).Post("/v3/login", s.apiLoginPost)
	r.With(middleware.Bind[model.RequestRegister]()).Post("/v3/register", s.apiRegister)
	r.With(middleware.AuthRequired(s.db)).Post("/v3/logout", s.apiLogout)
	r.With(middleware.AuthRequired(s.db)).Post("/v3/logout/all", s.apiLogoutAll)
	r.With(middleware.AuthOptional(s.db), middleware.Bind[model.RequestPassword]()).Post("/v3/account/password", s.apiChangePassword)
	r.With(middleware.AuthRequired(s.db)).Get("/v3/account/whoami", s.apiWhoami)

	r.Get("/v3/profile/{userId}", s.apiProfileDefault())
	r.Get("/v3/profile/{userId}/avatar_url", s.apiProfileAvatarURLGet())
	r.With(middleware.AuthRequired(s.db), middleware.Bind[model.RequestProfileUpdate]()).Put("/v3/profile/{userId}/avatar_url", s.apiProfileAvatarUrlPut())
	r.Get("/v3/profile/{userId}/displayname", s.apiProfileDisplayNameGet())
	r.With(middleware.AuthRequired(s.db), middleware.Bind[model.RequestProfileUpdate]()).Put("/v3/profile/{userId}/displayname", s.apiProfileDisplayNamePut())
	r.With(middleware.AuthRequired(s.db), middleware.Bind[model.RequestUserDirectorySearch]()).Post("/v3/user_directory/search", s.apiUserDirectorySearch)
	return r
}
