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
	r.With(middleware.RateLimiter()).Get("/v3/register/available", s.apiRegisterAvailable)
	r.With(middleware.RateLimiter()).Get("/v3/login", s.apiLoginGet)
	r.With(middleware.RateLimiter(), middleware.Bind[model.RequestLogin]()).Post("/v3/login", s.apiLoginPost)
	r.With(middleware.RateLimiter(), middleware.Bind[model.RequestRegister]()).Post("/v3/register", s.apiRegister)
	r.With(middleware.AuthRequired(s.db)).Post("/v3/logout", s.apiLogout)
	r.With(middleware.AuthRequired(s.db)).Post("/v3/logout/all", s.apiLogoutAll)
	r.With(middleware.RateLimiter(), middleware.AuthOptional(s.db), middleware.Bind[model.RequestPassword]()).Post("/v3/account/password", s.apiChangePassword)
	r.With(middleware.RateLimiter(), middleware.AuthRequired(s.db)).Get("/v3/account/whoami", s.apiWhoami)

	r.Get("/v3/profile/{userId}", s.apiProfileDefault())
	r.Get("/v3/profile/{userId}/avatar_url", s.apiProfileAvatarURLGet())
	r.With(middleware.RateLimiter(), middleware.AuthRequired(s.db), middleware.Bind[model.RequestProfileUpdate]()).Put("/v3/profile/{userId}/avatar_url", s.apiProfileAvatarUrlPut())
	r.Get("/v3/profile/{userId}/displayname", s.apiProfileDisplayNameGet())
	r.With(middleware.RateLimiter(), middleware.AuthRequired(s.db), middleware.Bind[model.RequestProfileUpdate]()).Put("/v3/profile/{userId}/displayname", s.apiProfileDisplayNamePut())
	r.With(middleware.RateLimiter(), middleware.AuthRequired(s.db), middleware.Bind[model.RequestUserDirectorySearch]()).Post("/v3/user_directory/search", s.apiUserDirectorySearch)

	r.With(middleware.RateLimiter(), middleware.AuthRequired(s.db)).Get("/v3/capabilities", s.apiCapabilitiesNegotiation)

	// r.With(middleware.AuthRequired(s.db), middleware.Bind[model.RequestCreateRoom]()).Post("/v3/createRoom")
	r.Get("/v3/directory/room/{roomAlias}", s.apiDirectoryRoomGet)
	r.With(middleware.AuthRequired(s.db), middleware.Bind[model.RequestDirectoryRoomPut]()).Put("/v3/directory/room/{roomAlias}", s.apiDirectoryRoomPut)
	r.With(middleware.AuthRequired(s.db)).Delete("/v3/directory/room/{roomAlias}", s.apiDirectoryRoomDelete)
	r.With(middleware.RateLimiter(), middleware.AuthRequired(s.db)).Get("/v3/rooms/{roomId}/aliases", s.apiRoomsAliases)

	// Events
	r.With(middleware.AuthRequired(s.db)).Get("/v3/sync", s.sync)
	r.With(middleware.AuthRequired(s.db)).Get("/v3/rooms/{roomId}/event/{eventId}", s.apiGetEvent)
	r.With(middleware.AuthRequired(s.db)).Get("/v3/rooms/{roomId}/joined_members", s.apiGetJoinedMembers)
	r.With(middleware.AuthRequired(s.db)).Get("/v3/rooms/{roomId}/members", s.apiGetMembers)
	// r.With(middleware.AuthRequired(s.db)).Get("/v3/rooms/{roomID}/state", s.)
	// r.With(middleware.AuthRequired(s.db)).Get("/v3/rooms/{roomID}/state/{eventType}/{stateKey}", s.sync)
	r.With(middleware.AuthRequired(s.db)).Get("/v3/rooms/{roomId}/messages", s.apiGetMessage)
	// r.With(middleware.RateLimiter(), middleware.AuthRequired(s.db)).Get("/v3/rooms/{roomID}/timestamp_to_event", s.sync)

	// Send Event
	// r.With(middleware.AuthRequired(s.db)).Put("/v3/rooms/{roomID}/state/{eventType}/{stateKey}", s.sync)
	r.With(middleware.AuthRequired(s.db), middleware.Bind[model.RequestSendMessage]()).Put("/v3/rooms/{roomID}/send/{eventType}/{txnID}", s.apiSend)
	// r.With(middleware.AuthRequired(s.db)).Put("/v3/rooms/{roomID}/redact/{eventID}/{txnID}", s.sync)
	return r
}
