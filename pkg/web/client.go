package web

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"hotaru.hana.im/server/pkg/crypto"
	"hotaru.hana.im/server/pkg/storage"
	"hotaru.hana.im/server/pkg/web/middleware"
	"hotaru.hana.im/server/pkg/web/model"
)

func clientVersionsHandler(w http.ResponseWriter, _ *http.Request) {
	resp := &model.ResponseClientVersions{
		Versions: []string{"v1.11"},
	}
	middleware.RenderJSON(w, resp)
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

	var accessToken string
	if err := s.db.IsDeviceExist(form.Identifier.User, form.DeviceID); err != nil {
		if !errors.Is(err, storage.ErrDeviceNotExist) {
			middleware.ErrorUnknown(w)
			return
		}
		if form.DeviceID == "" {
			form.DeviceID, _ = crypto.GenerateDeviceID()
		}
		accessToken, err = s.db.CreateDevice(form.Identifier.User, form.DeviceID, form.InitialDeviceDisplayName)
		if err != nil {
			middleware.ErrorUnknown(w)
			return
		}
	} else {
		accessToken, err = s.db.UpdateAccessToken(form.Identifier.User, form.DeviceID, form.InitialDeviceDisplayName)
		if err != nil {
			middleware.ErrorUnknown(w)
			return
		}
	}

	resp := &model.ResponseLoginPost{
		UserID:      s.toUserID(form.Identifier.User),
		AccessToken: accessToken,
		DeviceID:    form.DeviceID,
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

	if form.DeviceID == "" {
		form.DeviceID, _ = crypto.GenerateDeviceID()
	}
	accessToken, err := s.db.CreateDevice(form.Username, form.DeviceID, form.InitialDeviceDisplayName)
	if err != nil {
		middleware.ErrorUnknown(w)
		return
	}

	resp := &model.ResponseRegister{
		UserID:      s.toUserID(form.Username),
		AccessToken: accessToken,
		DeviceID:    form.DeviceID,
	}
	middleware.RenderJSON(w, resp)
}

func (s *Server) apiLogout(w http.ResponseWriter, r *http.Request) {
	device := middleware.GetDevice(r)
	err := s.db.DeleteDevice(device.Localpart, device.DeviceID)
	if err != nil {
		middleware.ErrorUnknownMsg(w, "unknown error")
		return
	}
	middleware.RenderJSON(w, &struct{}{})
}

func (s *Server) apiLogoutAll(w http.ResponseWriter, r *http.Request) {
	device := middleware.GetDevice(r)
	err := s.db.DeleteDeviceByLocalpart(device.Localpart)
	if err != nil {
		middleware.ErrorUnknownMsg(w, "unknown error")
		return
	}
	middleware.RenderJSON(w, &struct{}{})
}

func (s *Server) apiChangePassword(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		// TODO: use user-interactive authentication api
		middleware.ErrorForbidden(w)
	}
	form := middleware.GetObject(r).(*model.RequestPassword)
	account := middleware.GetAccount(r)

	// TODO: soft logout
	if err := s.db.ChangeAccountPassword(account, form.NewPassword); err != nil {
		middleware.ErrorUnknownMsg(w, "unknown error")
		return
	}

	middleware.RenderJSON(w, &struct{}{})
}

func (s *Server) apiWhoami(w http.ResponseWriter, r *http.Request) {
	account := middleware.GetAccount(r)
	device := middleware.GetDevice(r)

	resp := &model.ResponseWhoami{
		DeviceID: device.DeviceID,
		IsGuest:  false,
		UserID:   s.toUserID(account.Localpart),
	}
	middleware.RenderJSON(w, resp)
}

func (s *Server) apiProfileDefault() func(w http.ResponseWriter, r *http.Request) {
	return s.apiProfile(true, true)
}

func (s *Server) apiProfileDisplayNameGet() func(w http.ResponseWriter, r *http.Request) {
	return s.apiProfile(true, false)
}

func (s *Server) apiProfileAvatarURLGet() func(w http.ResponseWriter, r *http.Request) {
	return s.apiProfile(false, true)
}

func (s *Server) apiProfile(enableDisplayName, enableAvatarURL bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userId")
		localpart, domain, err := parseUserID(userID)
		if err != nil {
			middleware.ErrorNotFoundMsg(w, "invalid user id")
			return
		}

		// TODO: check user info of other server
		if domain != s.domain {
			middleware.ErrorForbidden(w)
		}

		profile, err := s.db.GetProfile(localpart)
		if err != nil {
			if errors.Is(err, storage.ErrProfileNotExist) {
				middleware.ErrorNotFoundMsg(w, "user profile not found")
			} else {
				middleware.ErrorUnknownMsg(w, "unknown error")
			}
			return
		}

		if !enableDisplayName {
			profile.DisplayName = ""
		}
		if !enableAvatarURL {
			profile.AvatarUrl = ""
		}

		resp := &model.ResponseProfile{
			DisplayName: profile.DisplayName,
			AvatarURL:   profile.AvatarUrl,
		}
		middleware.RenderJSON(w, resp)
	}
}

func (s *Server) toUserID(localpart string) string {
	return "@" + localpart + ":" + s.domain
}

var userIDParser = regexp.MustCompile(`^@([0-9a-z_\-+=./]+):(.+)$`)
var ErrInvalidUserID = errors.New("invalid user id")

func parseUserID(userID string) (localpart, domain string, err error) {
	matches := userIDParser.FindStringSubmatch(userID)
	if len(matches) == 0 {
		return "", "", ErrInvalidUserID
	}
	return matches[1], matches[2], nil
}
