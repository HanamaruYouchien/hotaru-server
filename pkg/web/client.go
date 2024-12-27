package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"hotaru.hana.im/server/pkg/crypto"
	"hotaru.hana.im/server/pkg/storage"
	"hotaru.hana.im/server/pkg/web/middleware"
	"hotaru.hana.im/server/pkg/web/model"
)

func (s *Server) toUserID(localpart string) string {
	return fmt.Sprintf("@%s:%s", localpart, s.domain)
}

func (s *Server) toRoomAlias(alias string) string {
	return fmt.Sprintf("#%s:%s", alias, s.domain)
}

func (s *Server) toRoomID(roomID string) string {
	return fmt.Sprintf("!%s:%s", roomID, s.domain)
}

func (s *Server) toEventID(eventID string) string {
	return fmt.Sprintf("$%s:%s", eventID, s.domain)
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

var aliasParser = regexp.MustCompile(`^#(.+):(.+)$`)
var ErrInvalidRoomAlias = errors.New("invalid room alias")

func parseAlias(roomAlias string) (alias, domain string, err error) {
	matches := aliasParser.FindStringSubmatch(roomAlias)
	if len(matches) == 0 {
		return "", "", ErrInvalidUserID
	}
	return matches[1], matches[2], nil
}

var roomIDParser = regexp.MustCompile(`^!(.+):(.+)$`)
var ErrInvalidRoomID = errors.New("invalid room id")

func parseRoomID(room string) (roomID, domain string, err error) {
	matches := roomIDParser.FindStringSubmatch(room)
	if len(matches) == 0 {
		return "", "", ErrInvalidRoomID
	}
	return matches[1], matches[2], nil
}

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
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
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
		middleware.ErrorUnknown(w, http.StatusBadRequest)
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

	// TODO: try to update first?
	var accessToken string
	if err := s.db.IsDeviceExist(form.Identifier.User, form.DeviceID); err != nil {
		if !errors.Is(err, storage.ErrDeviceNotExist) {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
			return
		}
		if form.DeviceID == "" {
			form.DeviceID, _ = crypto.GenerateDeviceID()
		}
		accessToken, err = s.db.CreateDevice(form.Identifier.User, form.DeviceID, form.InitialDeviceDisplayName)
		if err != nil {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
			return
		}
	} else {
		accessToken, err = s.db.UpdateAccessToken(form.Identifier.User, form.DeviceID, form.InitialDeviceDisplayName)
		if err != nil {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
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
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}

	if form.DeviceID == "" {
		form.DeviceID, _ = crypto.GenerateDeviceID()
	}
	accessToken, err := s.db.CreateDevice(form.Username, form.DeviceID, form.InitialDeviceDisplayName)
	if err != nil {
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
		return
	}

	s.db.CreateProfile(form.Username, form.Username, "")

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
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
		return
	}
	middleware.RenderJSON(w, &struct{}{})
}

func (s *Server) apiLogoutAll(w http.ResponseWriter, r *http.Request) {
	device := middleware.GetDevice(r)
	err := s.db.DeleteDeviceByLocalpart(device.Localpart)
	if err != nil {
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
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
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
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

		// TODO: get user profile from other server
		if domain != s.domain {
			middleware.ErrorForbidden(w)
			return
		}

		profile, err := s.db.GetProfile(localpart)
		if err != nil {
			if errors.Is(err, storage.ErrProfileNotExist) {
				middleware.ErrorNotFoundMsg(w, "user profile not found")
			} else {
				middleware.ErrorUnknown(w, http.StatusInternalServerError)
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

func (s *Server) apiProfileDisplayNamePut() func(w http.ResponseWriter, r *http.Request) {
	return s.apiProfileUpdate(true, false)
}

func (s *Server) apiProfileAvatarUrlPut() func(w http.ResponseWriter, r *http.Request) {
	return s.apiProfileUpdate(false, true)
}

func (s *Server) apiProfileUpdate(enableDisplayName, enableAvatarURL bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userId")
		localpart, domain, err := parseUserID(userID)
		if err != nil {
			middleware.ErrorNotFoundMsg(w, "invalid user id")
			return
		}

		// TODO: get user profile from other server
		if domain != s.domain {
			middleware.ErrorForbidden(w)
			return
		}

		account := middleware.GetAccount(r)
		if account.Localpart != localpart {
			middleware.ErrorForbiddenMsg(w, "no permission")
			return
		}

		form := middleware.GetObject(r).(*model.RequestProfileUpdate)
		if !enableDisplayName {
			form.DisplayName = ""
		}
		if !enableAvatarURL {
			form.AvatarUrl = ""
		}

		if enableDisplayName {
			err = s.db.UpdateProfileDisplayName(localpart, form.DisplayName)
		}
		if enableAvatarURL {
			err = s.db.UpdateProfileAvatarUrl(localpart, form.AvatarUrl)
		}
		if err == nil {
			middleware.RenderJSON(w, &struct{}{})
			return
		}

		if err := s.db.IsProfileExist(localpart); err != nil {
			if !errors.Is(err, storage.ErrProfileNotExist) {
				middleware.ErrorUnknown(w, http.StatusInternalServerError)
				return
			}
			if err := s.db.CreateProfile(localpart, form.DisplayName, form.AvatarUrl); err != nil {
				middleware.ErrorUnknown(w, http.StatusInternalServerError)
				return
			}
		} else {
			if enableDisplayName {
				if err := s.db.UpdateProfileDisplayName(localpart, form.DisplayName); err != nil {
					middleware.ErrorUnknown(w, http.StatusInternalServerError)
					return
				}
			}
			if enableAvatarURL {
				if err := s.db.UpdateProfileAvatarUrl(localpart, form.AvatarUrl); err != nil {
					middleware.ErrorUnknown(w, http.StatusInternalServerError)
					return
				}
			}
		}

		middleware.RenderJSON(w, &struct{}{})
	}
}

func (s *Server) apiUserDirectorySearch(w http.ResponseWriter, r *http.Request) {
	form := middleware.GetObject(r).(*model.RequestUserDirectorySearch)
	if form.SearchTerm == "" {
		middleware.RenderJSON(w, &model.ResponseUserDirectorySearch{
			Results: []model.User{},
		})
		return
	}
	if form.Limit == 0 { // default value
		form.Limit = 10
	}

	profiles, err := s.db.SearchProfile(form.SearchTerm, form.Limit)
	if err != nil {
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
		return
	}

	resp := &model.ResponseUserDirectorySearch{
		Limited: len(profiles) == form.Limit,
		Results: make([]model.User, 0, len(profiles)),
	}
	for _, v := range profiles {
		resp.Results = append(resp.Results, model.User{
			UserID:      s.toUserID(v.Localpart),
			AvatarURL:   v.AvatarUrl,
			DisplayName: v.DisplayName,
		})
	}

	middleware.RenderJSON(w, resp)
}

func (s *Server) apiCapabilitiesNegotiation(w http.ResponseWriter, r *http.Request) {
	response := model.CapabilitiesResponse{}

	response.Capabilities.ThreePIDChanges = model.BooleanCapability{Enabled: false}
	response.Capabilities.ChangePassword = model.BooleanCapability{Enabled: true}
	response.Capabilities.GetLoginToken = model.BooleanCapability{Enabled: false}
	response.Capabilities.RoomVersions = model.RoomVersionsCapability{
		Default:   "v10",
		Available: []string{"v10"},
	}
	response.Capabilities.SetAvatarURL = model.BooleanCapability{Enabled: false}
	response.Capabilities.SetDisplayName = model.BooleanCapability{Enabled: true}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
		return
	}
}

func (s *Server) apiCreateRoom(w http.ResponseWriter, r *http.Request) {
	form := middleware.GetObject(r).(*model.RequestCreateRoom)
	account := middleware.GetAccount(r)

	roomID, err := s.db.CreateRoom(form.Name, form.Topic, form.Visibility, account.Localpart, form.RoomAliasName)
	if err != nil {
		if errors.Is(err, storage.ErrRoomAliasInUsed) {
			middleware.ErrorRoomInUse(w)
		} else {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}

	resp := &model.ResponseCreateRoom{RoomID: s.toRoomID(roomID)}
	middleware.RenderJSON(w, resp)
}

func (s *Server) apiDirectoryRoomGet(w http.ResponseWriter, r *http.Request) {
	roomAlias := chi.URLParam(r, "roomAlias")
	alias, domain, err := parseAlias(roomAlias)
	if err != nil || domain != s.domain {
		middleware.ErrorInvalidParamMsg(w, "room alias invalid")
		return
	}

	roomID, err := s.db.GetRoomIDByRoomAlias(alias)
	if err != nil {
		if errors.Is(err, storage.ErrRoomAliasNotExist) {
			middleware.ErrorNotFoundMsg(w, "room alias not found")
		} else {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
			return
		}
	}

	resp := &model.ResponseDirectoryRoomGet{
		RoomID:  s.toRoomID(roomID),
		Servers: []string{s.domain},
	}
	middleware.RenderJSON(w, resp)
}

func (s *Server) apiDirectoryRoomPut(w http.ResponseWriter, r *http.Request) {
	roomAlias := chi.URLParam(r, "roomAlias")
	alias, domain, err := parseAlias(roomAlias)
	if err != nil || domain != s.domain {
		middleware.ErrorInvalidParamMsg(w, "room alias invalid")
		return
	}

	form := middleware.GetObject(r).(*model.RequestDirectoryRoomPut)
	account := middleware.GetAccount(r)

	roomID, domain, err := parseRoomID(form.RoomID)
	if err != nil || domain != s.domain {
		middleware.ErrorInvalidParamMsg(w, "room id invalid")
		return
	}

	if err := s.db.CheckPermission(roomID, account.Localpart, 100); err != nil {
		if errors.Is(err, storage.ErrAccountNotInRoom) || errors.Is(err, storage.ErrNoPermission) {
			middleware.ErrorForbidden(w)
		} else {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}

	if err := s.db.IsRoomAliasExist(alias); err == nil {
		middleware.ErrorUnknownMsg(w, http.StatusConflict, "room alias already exists")
		return
	} else if !errors.Is(err, storage.ErrRoomAliasNotExist) {
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
		return
	}

	if err := s.db.CreateRoomAlias(alias, roomID); err != nil {
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
		return
	}
	middleware.RenderJSON(w, &struct{}{})
}

func (s *Server) apiDirectoryRoomDelete(w http.ResponseWriter, r *http.Request) {
	roomAlias := chi.URLParam(r, "roomAlias")
	alias, domain, err := parseAlias(roomAlias)
	if err != nil || domain != s.domain {
		middleware.ErrorInvalidParamMsg(w, "room alias invalid")
		return
	}

	roomID, err := s.db.GetRoomIDByRoomAlias(alias)
	if err != nil {
		if errors.Is(err, storage.ErrRoomAliasNotExist) {
			middleware.ErrorNotFoundMsg(w, "room alias not found")
		} else {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}

	account := middleware.GetAccount(r)
	if err := s.db.CheckPermission(roomID, account.Localpart, 100); err != nil {
		if errors.Is(err, storage.ErrAccountNotInRoom) || errors.Is(err, storage.ErrNoPermission) {
			middleware.ErrorForbidden(w)
		} else {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}

	if err := s.db.DeleteRoomAlias(alias); err != nil {
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
		return
	}
	middleware.RenderJSON(w, &struct{}{})
}

func (s *Server) apiRoomsAliases(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	roomID, domain, err := parseRoomID(roomID)
	if err != nil || domain != s.domain {
		middleware.ErrorInvalidParamMsg(w, "room id invalid")
		return
	}

	account := middleware.GetAccount(r)
	if err := s.db.CheckPermission(roomID, account.Localpart, 100); err != nil {
		if errors.Is(err, storage.ErrAccountNotInRoom) || errors.Is(err, storage.ErrNoPermission) {
			middleware.ErrorForbidden(w)
		} else {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}

	aliases, err := s.db.GetRoomAliasesByRoomID(roomID)
	if err != nil {
		if errors.Is(err, storage.ErrEmptyRoomID) {
			middleware.ErrorForbiddenMsg(w, "roomID is empty")
		} else {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}

	for k, v := range aliases {
		aliases[k] = s.toRoomAlias(v)
	}
	resp := &model.ResponseRoomsAliases{Aliases: aliases}
	middleware.RenderJSON(w, resp)
}

func (s *Server) apiJoinedRooms(w http.ResponseWriter, r *http.Request) {
	account := middleware.GetAccount(r)
	joinedRooms, err := s.db.GetJoinedRooms(account.Localpart)
	if err != nil {
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
		return
	}

	for k, v := range joinedRooms {
		joinedRooms[k] = s.toRoomID(v)
	}

	resp := &model.ResponseJoinedRooms{JoinedRooms: joinedRooms}
	middleware.RenderJSON(w, resp)
}

func (s *Server) apiRoomsInvite(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	roomID, domain, err := parseRoomID(roomID)
	if err != nil || domain != s.domain {
		middleware.ErrorInvalidParamMsg(w, "room id invalid")
		return
	}
	account := middleware.GetAccount(r)
	form := middleware.GetObject(r).(*model.RequestRoomsInvite)
	userID, domain, err := parseUserID(form.UserID)
	if err != nil {
		middleware.ErrorUnknown(w, http.StatusBadRequest)
		return
	}

	// TODO: reason

	if err := s.db.RoomInvite(roomID, account.Localpart, userID); err != nil {
		switch {
		case errors.Is(err, storage.ErrEmptyLocalpart):
			fallthrough
		case errors.Is(err, storage.ErrEmptyRoomID):
			middleware.ErrorUnknown(w, http.StatusBadRequest)
		case errors.Is(err, storage.ErrNoPermission):
			middleware.ErrorForbidden(w)
		default:
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}
	middleware.RenderJSON(w, &struct{}{})
}

func (s *Server) apiRoomsJoin(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	roomID, domain, err := parseRoomID(roomID)
	if err != nil || domain != s.domain {
		middleware.ErrorForbiddenMsg(w, "room id invalid")
		return
	}

	account := middleware.GetAccount(r)
	form := middleware.GetObject(r).(*model.RequestRoomsJoin)
	_ = form

	// TODO: reason

	if err := s.db.RoomJoin(roomID, account.Localpart); err != nil {
		switch {
		case errors.Is(err, storage.ErrEmptyLocalpart):
			fallthrough
		case errors.Is(err, storage.ErrEmptyRoomID):
			fallthrough
		case errors.Is(err, storage.ErrNoPermission):
			middleware.ErrorForbidden(w)
		default:
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}

	resp := &model.ResponseRoomsJoin{RoomID: s.toRoomID(roomID)}
	middleware.RenderJSON(w, resp)
}

func (s *Server) apiJoin(w http.ResponseWriter, r *http.Request) {
	roomIdOrAlias := chi.URLParam(r, "roomIdOrAlias")
	alias, domain, err := parseAlias(roomIdOrAlias)
	var roomID = ""
	if err == nil {
		if domain != s.domain {
			middleware.ErrorForbiddenMsg(w, "room id or alias invalid")
			return
		}

		roomID, err = s.db.GetRoomIDByRoomAlias(alias)
		if err != nil {
			if errors.Is(err, storage.ErrRoomAliasNotExist) {
				middleware.ErrorForbiddenMsg(w, "room alias not found")
			} else {
				middleware.ErrorUnknown(w, http.StatusInternalServerError)
				return
			}
		}
	} else {
		roomID, domain, err = parseRoomID(roomID)
		if err != nil || domain != s.domain {
			middleware.ErrorForbiddenMsg(w, "room id invalid")
			return
		}
	}

	account := middleware.GetAccount(r)
	form := middleware.GetObject(r).(*model.RequestJoin)
	_ = form

	// TODO: reason

	if err := s.db.RoomJoin(roomID, account.Localpart); err != nil {
		switch {
		case errors.Is(err, storage.ErrEmptyLocalpart):
			fallthrough
		case errors.Is(err, storage.ErrEmptyRoomID):
			fallthrough
		case errors.Is(err, storage.ErrNoPermission):
			middleware.ErrorForbidden(w)
		default:
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}

	resp := &model.ResponseJoinRoom{RoomID: s.toRoomID(roomID)}
	middleware.RenderJSON(w, resp)
}

func (s *Server) apiRoomsLeave(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	roomID, domain, err := parseRoomID(roomID)
	if err != nil || domain != s.domain {
		middleware.ErrorInvalidParamMsg(w, "room id invalid")
		return
	}

	account := middleware.GetAccount(r)
	form := middleware.GetObject(r).(*model.RequestRoomsLeave)
	_ = form

	if err := s.db.RoomLeave(roomID, account.Localpart); err != nil {
		switch {
		case errors.Is(err, storage.ErrEmptyLocalpart):
			fallthrough
		case errors.Is(err, storage.ErrEmptyRoomID):
			middleware.ErrorInvalidParam(w)
		case errors.Is(err, storage.ErrNoPermission):
			middleware.ErrorForbidden(w)
		default:
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}
	middleware.RenderJSON(w, &struct{}{})
}

func (s *Server) apiRoomsKick(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	roomID, domain, err := parseRoomID(roomID)
	if err != nil || domain != s.domain {
		middleware.ErrorInvalidParamMsg(w, "room id invalid")
		return
	}

	account := middleware.GetAccount(r)
	form := middleware.GetObject(r).(*model.RequestRoomsKick)
	userID, domain, err := parseUserID(form.UserID)
	if err != nil {
		middleware.ErrorUnknown(w, http.StatusBadRequest)
		return
	}

	if err := s.db.RoomKick(roomID, account.Localpart, userID); err != nil {
		switch {
		case errors.Is(err, storage.ErrEmptyLocalpart):
			fallthrough
		case errors.Is(err, storage.ErrEmptyRoomID):
			middleware.ErrorInvalidParam(w)
		case errors.Is(err, storage.ErrNoPermission):
			middleware.ErrorForbidden(w)
		default:
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}
	middleware.RenderJSON(w, &struct{}{})
}

func (s *Server) apiRoomsForget(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	roomID, domain, err := parseRoomID(roomID)
	if err != nil || domain != s.domain {
		middleware.ErrorInvalidParamMsg(w, "room id invalid")
		return
	}

	account := middleware.GetAccount(r)

	if err := s.db.RoomForget(roomID, account.Localpart); err != nil {
		switch {
		case errors.Is(err, storage.ErrEmptyLocalpart):
			fallthrough
		case errors.Is(err, storage.ErrEmptyRoomID):
			middleware.ErrorInvalidParam(w)
		case errors.Is(err, storage.ErrNoPermission):
			middleware.ErrorUnknownMsg(w, http.StatusBadRequest, "user in room")
		default:
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}
	middleware.RenderJSON(w, &struct{}{})
}

func (s *Server) apiRoomsBan(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	roomID, domain, err := parseRoomID(roomID)
	if err != nil || domain != s.domain {
		middleware.ErrorInvalidParamMsg(w, "room id invalid")
		return
	}

	account := middleware.GetAccount(r)
	form := middleware.GetObject(r).(*model.RequestRoomsBan)
	userID, domain, err := parseUserID(form.UserID)
	if err != nil {
		middleware.ErrorUnknown(w, http.StatusBadRequest)
		return
	}

	if err := s.db.RoomBan(roomID, account.Localpart, userID); err != nil {
		switch {
		case errors.Is(err, storage.ErrEmptyLocalpart):
			fallthrough
		case errors.Is(err, storage.ErrEmptyRoomID):
			middleware.ErrorInvalidParam(w)
		case errors.Is(err, storage.ErrNoPermission):
			middleware.ErrorForbidden(w)
		default:
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}
	middleware.RenderJSON(w, &struct{}{})
}

func (s *Server) apiRoomsUnban(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	roomID, domain, err := parseRoomID(roomID)
	if err != nil || domain != s.domain {
		middleware.ErrorInvalidParamMsg(w, "room id invalid")
		return
	}

	account := middleware.GetAccount(r)
	form := middleware.GetObject(r).(*model.RequestRoomsUnban)
	userID, domain, err := parseUserID(form.UserID)
	if err != nil {
		middleware.ErrorUnknown(w, http.StatusBadRequest)
		return
	}

	if err := s.db.RoomUnban(roomID, account.Localpart, userID); err != nil {
		switch {
		case errors.Is(err, storage.ErrEmptyLocalpart):
			fallthrough
		case errors.Is(err, storage.ErrEmptyRoomID):
			middleware.ErrorInvalidParam(w)
		case errors.Is(err, storage.ErrNoPermission):
			middleware.ErrorForbidden(w)
		default:
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}
	middleware.RenderJSON(w, &struct{}{})
}

func (s *Server) apiDirectoryListRoomGet(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	roomID, domain, err := parseRoomID(roomID)
	if err != nil || domain != s.domain {
		middleware.ErrorInvalidParamMsg(w, "room id invalid")
		return
	}

	room, err := s.db.GetRoom(roomID)
	if err != nil {
		if errors.Is(err, storage.ErrRoomNotExist) {
			middleware.ErrorNotFound(w)
		} else {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}

	resp := &model.ResponseDirectoryListRoomGet{
		Visibility: room.Visibility,
	}
	middleware.RenderJSON(w, resp)
}

func (s *Server) apiDirectoryListRoomPut(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	roomID, domain, err := parseRoomID(roomID)
	if err != nil || domain != s.domain {
		middleware.ErrorInvalidParamMsg(w, "room id invalid")
		return
	}

	account := middleware.GetAccount(r)
	form := middleware.GetObject(r).(*model.RequestDirectoryListRoomPut)

	if err := s.db.CheckPermission(roomID, account.Localpart, 100); err != nil {
		if errors.Is(err, storage.ErrAccountNotInRoom) || errors.Is(err, storage.ErrNoPermission) {
			middleware.ErrorForbidden(w)
		} else {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}

	if err := s.db.IsRoomExist(roomID); err != nil {
		if errors.Is(err, storage.ErrRoomNotExist) {
			middleware.ErrorNotFound(w)
		} else {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}
	if err := s.db.UpdateRoomVisibility(roomID, form.Visibility); err != nil {
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
		return
	}
	middleware.RenderJSON(w, &struct{}{})
}

func (s *Server) apiPublicRoomsGet(w http.ResponseWriter, _ *http.Request) {
	// TODO: term of search
	rooms, err := s.db.GetPublicRooms()
	if err != nil {
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
		return
	}
	resp := &model.ResponsePublicRoomsGet{
		Chunk: make([]model.PublicRoomsChunk, 0, len(rooms)),
	}

	for _, v := range rooms {
		resp.Chunk = append(resp.Chunk, model.PublicRoomsChunk{
			RoomID:           s.toRoomID(v.RoomId),
			Name:             v.Name,
			NumJoinedMembers: 0,
			WorldReadable:    false,
			GuestCanJoin:     false,
		})
	}
	middleware.RenderJSON(w, resp)
}

func (s *Server) apiPublicRoomsPost(w http.ResponseWriter, _ *http.Request) {
	// TODO: term of search
	rooms, err := s.db.GetPublicRooms()
	if err != nil {
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
		return
	}
	resp := &model.ResponsePublicRoomsPost{
		Chunk: make([]model.PublicRoomsChunk, 0, len(rooms)),
	}

	for _, v := range rooms {
		resp.Chunk = append(resp.Chunk, model.PublicRoomsChunk{
			RoomID:           s.toRoomID(v.RoomId),
			Name:             v.Name,
			NumJoinedMembers: 0,
			WorldReadable:    false,
			GuestCanJoin:     false,
		})
	}
	middleware.RenderJSON(w, resp)
}

func (s *Server) formatEvent(inputEvent storage.Event) storage.Event {
	event := inputEvent
	event.EventId = s.toEventID(inputEvent.EventId)
	event.RoomId = s.toRoomID(inputEvent.RoomId)
	event.Sender = s.toUserID(inputEvent.Sender)
	event.OriginServerTs = event.CreatedAt.UnixMicro()
	return event
}

// Get Event
func (s *Server) sync(w http.ResponseWriter, r *http.Request) {
	account := middleware.GetAccount(r)
	_ = account
	// TODO do it
}

func (s *Server) apiGetEvent(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	eventID := chi.URLParam(r, "eventId")
	if err := s.db.IsAccountInRoom(roomID, middleware.GetAccount(r).Localpart); err != nil {
		if errors.Is(err, storage.ErrAccountNotInRoom) {
			middleware.ErrorForbiddenMsg(w, "You aren’t a member of the room and weren’t previously a member of the room.")
		} else {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}

	event, err := s.db.GetEvent(roomID, eventID)
	if err != nil {
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
		return
	}

	res := model.EventResponse{}
	res.Event = s.formatEvent(*event)

	middleware.RenderJSON(w, res.Event)
}

func (s *Server) apiGetJoinedMembers(w http.ResponseWriter, r *http.Request) {
	roomID, domain, err := parseRoomID(chi.URLParam(r, "roomId"))
	if err != nil || domain != s.domain {
		middleware.ErrorInvalidParamMsg(w, "room id invalid")
		return
	}
	if err := s.db.IsAccountInRoom(roomID, middleware.GetAccount(r).Localpart); err != nil {
		if errors.Is(err, storage.ErrAccountNotInRoom) {
			middleware.ErrorForbiddenMsg(w, "You aren’t a member of the room and weren’t previously a member of the room.")
		} else {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}

	res, err := s.db.GetJoinedMembers(roomID)
	if err != nil {
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
		return
	}

	JoinedMembers := model.JoinedMembersResponse{
		Joined: make(map[string]model.Member),
	}

	for _, username := range res {
		JoinedMembers.Joined[s.toUserID(username)] = model.Member{
			AvatarURL:   nil,
			DisplayName: username,
		}
	}

	middleware.RenderJSON(w, JoinedMembers)
}

func (s *Server) apiGetMembers(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	if err := s.db.IsAccountInRoom(roomID, middleware.GetAccount(r).Localpart); err != nil {
		if errors.Is(err, storage.ErrAccountNotInRoom) {
			middleware.ErrorForbiddenMsg(w, "You aren’t a member of the room and weren’t previously a member of the room.")
		} else {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}

	res, err := s.db.GetMembers(roomID)
	if err != nil {
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
		return
	}

	Members := model.MembersResponse{}
	for _, event := range res {
		Members.Events = append(Members.Events, s.formatEvent(event))
	}

	middleware.RenderJSON(w, Members)
}

func (s *Server) apiGetMessage(w http.ResponseWriter, r *http.Request) {
	roomID, domain, err := parseRoomID(chi.URLParam(r, "roomId"))
	if err != nil || domain != s.domain {
		middleware.ErrorInvalidParamMsg(w, "room id invalid")
		return
	}
	if err := s.db.IsAccountInRoom(roomID, middleware.GetAccount(r).Localpart); err != nil {
		if errors.Is(err, storage.ErrAccountNotInRoom) {
			middleware.ErrorForbiddenMsg(w, "You aren’t a member of the room and weren’t previously a member of the room.")
		} else {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}

	request := model.RequestGetMessage{}
	request.Dir = r.URL.Query().Get("dir")
	request.Filter = r.URL.Query().Get("filter")
	request.Limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	request.From, _ = time.Parse(time.RFC3339, r.URL.Query().Get("from"))

	toParam := r.URL.Query().Get("to")
	if toParam != "" {
		request.To, _ = time.Parse(time.RFC3339, toParam)
	} else {
		request.To = time.Now()
	}

	events, err := s.db.GetMessages(roomID, request.Limit, request.Dir, request.From, request.To)
	if err != nil {
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
		return
	}
	res := model.MessagesResponse{}
	res.Start = request.From.Format(time.RFC3339)
	res.End = request.To.Format(time.RFC3339)
	for _, event := range events {
		res.Events = append(res.Events, s.formatEvent(event))
	}
	middleware.RenderJSON(w, res)
}

// Send Event
func (s *Server) apiSend(w http.ResponseWriter, r *http.Request) {
	sender := middleware.GetAccount(r).Localpart
	roomID := chi.URLParam(r, "roomID")
	eventType := chi.URLParam(r, "eventType")
	txnID := chi.URLParam(r, "txnID")
	if err := s.db.IsAccountInRoom(roomID, sender); err != nil {
		if errors.Is(err, storage.ErrAccountNotInRoom) {
			middleware.ErrorForbiddenMsg(w, "You aren’t a member of the room and weren’t previously a member of the room.")
		} else {
			middleware.ErrorUnknown(w, http.StatusInternalServerError)
		}
		return
	}

	form := middleware.GetObject(r).(*model.RequestSendMessage)
	var eventID string
	var err error

	// TODO: check txnId
	switch form.MsgType {
	case storage.MessageTypeText, storage.MessageTypeEmote, storage.MessageTypeNotice:
		eventID, err = s.db.SendText(roomID, eventType, txnID, json.RawMessage(form.Body), sender)

	// case storage.MessageTypeImage:
	// 	requestBody := storage.EventRoomMessageImageContent{}
	// 	if errDecodeBody := json.Unmarshal(bodyBytes, &requestBody); errDecodeBody != nil {
	// 		middleware.ErrorUnknown(w, http.StatusBadRequest)
	// 		return
	// 	}
	// 	// eventID, err = s.db.SendImage(roomID, eventType, txnID, jsonString, "user01")

	// case storage.MessageTypeFile:
	// 	requestBody := storage.EventRoomMessageFileContent{}
	// 	if errDecodeBody := json.Unmarshal(bodyBytes, &requestBody); errDecodeBody != nil {
	// 		middleware.ErrorUnknown(w, http.StatusBadRequest)
	// 		return
	// 	}
	// 	// eventID, err = s.db.SendFile(roomID, eventType, txnID, jsonString, "user01")

	// case storage.MessageTypeAudio:
	// 	requestBody := storage.EventRoomMessageAudioContent{}
	// 	if errDecodeBody := json.Unmarshal(bodyBytes, &requestBody); errDecodeBody != nil {
	// 		middleware.ErrorUnknown(w, http.StatusBadRequest)
	// 		return
	// 	}
	// 	// eventID, err = s.db.SendAudio(roomID, eventType, txnID, jsonString, "user01")

	// case storage.MessageTypeVideo:
	// 	requestBody := storage.EventRoomMessageVideoContent{}
	// 	if errDecodeBody := json.Unmarshal(bodyBytes, &requestBody); errDecodeBody != nil {
	// 		middleware.ErrorUnknown(w, http.StatusBadRequest)
	// 		return
	// 	}
	// eventID, err = s.db.SendVideo(roomID, eventType, txnID, jsonString, "user01")

	default:
		middleware.ErrorUnknown(w, http.StatusBadRequest)
		return
	}

	if err != nil {
		middleware.ErrorUnknown(w, http.StatusInternalServerError)
		return
	}
	res := model.EventSentResponse{}
	res.EventID = eventID
	middleware.RenderJSON(w, res)
}
