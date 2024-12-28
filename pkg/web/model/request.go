package model

import (
	"time"

	"hotaru.hana.im/server/pkg/storage"
)

const (
	AuthenticationTypePassword           = "m.login.password"
	AuthenticationTypeRecaptcha          = "m.login.recaptcha"
	AuthenticationTypeSSO                = "m.login.sso"
	AuthenticationTypeIdentity           = "m.login.email.identity"
	AuthenticationTypeMsisdn             = "m.login.msisdn"
	AuthenticationTypeDummy              = "m.login.dummy"
	AuthenticationTypeRegistrationToken  = "m.login.registration_token"
	AuthenticationTypeApplicationService = "m.login.application_service"
	AuthenticationTypeToken              = "m.login.token"

	IdentifierTypeUser       = "m.id.user"
	IdentifierTypeThirdParty = "m.id.thirdparty"
	IdentifierTypePhone      = "m.id.phone"
)

type Identifier struct {
	Type    string `json:"type"`
	User    string `json:"user,omitempty"`
	Medium  string `json:"medium,omitempty"`
	Address string `json:"address,omitempty"`
	Country string `json:"country,omitempty"`
	Phone   string `json:"phone,omitempty"`
}

type RequestLogin struct {
	DeviceID                 string     `json:"device_id,omitempty"`
	Identifier               Identifier `json:"identifier"`
	InitialDeviceDisplayName string     `json:"initial_device_display_name,omitempty"`
	Password                 string     `json:"password,omitempty"`
	RefreshToken             bool       `json:"refresh_token,omitempty"`
	Token                    string     `json:"token,omitempty"`
	Type                     string     `json:"type"`
}

type RequestGetToken struct {
	Auth AuthenticationData `json:"auth"`
}

// Interactive Auth API
type AuthenticationData struct {
	Type          string          `json:"type"`
	Session       string          `json:"session"`
	Identifier    Identifier      `json:"identifier,omitempty"`
	ThreepidCreds Request3pidBind `json:"threepid_creds,omitempty"`
	Password      string          `json:"password,omitempty"`
	Response      string          `json:"response,omitempty"`
	Token         string          `json:"token,omitempty"`
	// TODO: Fallback
}

type RequestRefresh struct {
	RefreshToken string `json:"refresh_token"`
}

type RequestDeactivate struct {
	Auth     AuthenticationData `json:"auth,omitempty"`
	Erase    bool               `json:"erase,omitempty"`
	IDServer string             `json:"id_server,omitempty"`
}

type RequestPassword struct {
	Auth          AuthenticationData `json:"auth,omitempty"`
	LogoutDevices bool               `json:"logout_devices,omitempty"`
	NewPassword   string             `json:"new_password"`
}

// RequestPasswordEmailRequestToken, RequestRegisterEmailRequestToken, Request3pidEmailRequestToken
type RequestEmailRequestToken struct {
	ClientSecret  string `json:"client_secret"`
	Email         string `json:"email"`
	IDAccessToken string `json:"id_access_token,omitempty"`
	IDServer      string `json:"id_server,omitempty"`
	NextLink      string `json:"next_link,omitempty"`
	SendAttempt   int    `json:"send_attempt"`
}

// RequestPasswordMsisdnRequestToken, RequestRegisterMsisdnRequestToken, RequestMsisdnEmailRequestToken
type RequestMsisdnRequestToken struct {
	ClientSecret  string `json:"client_secret"`
	Country       string `json:"country"`
	IDAccessToken string `json:"id_access_token,omitempty"`
	IDServer      string `json:"id_server,omitempty"`
	NextLink      string `json:"next_link,omitempty"`
	PhoneNumber   string `json:"phone_number"`
	SendAttempt   int    `json:"send_attempt"`
}

type RequestRegister struct {
	Auth                     AuthenticationData `json:"auth,omitempty"`
	DeviceID                 string             `json:"device_id,omitempty"`
	InhibitLogin             bool               `json:"inhibit_login,omitempty"`
	InitialDeviceDisplayName string             `json:"initial_device_display_name,omitempty"`
	Password                 string             `json:"password,omitempty"`
	RefreshToken             bool               `json:"refresh_token,omitempty"`
	Username                 string             `json:"username"`
}

type Request3pidAdd struct {
	Auth         AuthenticationData `json:"auth,omitempty"`
	ClientSecret string             `json:"client_secret"`
	Sid          string             `json:"sid"`
}

type Request3pidBind struct {
	ClientSecret  string `json:"client_secret"`
	IDAccessToken string `json:"id_access_token"`
	IDServer      string `json:"id_server"`
	Sid           string `json:"sid"`
}

// Request3pidDelete, Request3pidUnbind
type Request3pidDelete struct {
	Address  string `json:"address"`
	IDServer string `json:"id_server,omitempty"`
	Medium   string `json:"medium"`
}
type Request3pidUnbind Request3pidDelete

type RequestProfileUpdate struct {
	DisplayName string `json:"displayname,omitempty"`
	AvatarUrl   string `json:"avatar_url,omitempty"`
}
type RequestProfileDisplayNamePut RequestProfileUpdate
type RequestProfileAvatarUrlPut RequestProfileUpdate

type RequestUserDirectorySearch struct {
	Limit      int    `json:"limit,omitempty"`
	SearchTerm string `json:"search_term"`
}

type RequestCreateRoom struct {
	CreationContent           storage.EventRoomCreateContent      `json:"creation_content,omitempty"`
	InitialState              []StateEvent                        `json:"initial_state,omitempty"`
	Invite                    []string                            `json:"invite,omitempty"`
	Invite3pid                []Invite3pid                        `json:"invite_3pid,omitempty"`
	IsDirect                  bool                                `json:"is_direct,omitempty"`
	Name                      string                              `json:"name,omitempty"`
	PowerLevelContentOverride storage.EventRoomPowerLevelsContent `json:"power_level_content_override,omitempty"`
	Preset                    string                              `json:"preset,omitempty"`
	RoomAliasName             string                              `json:"room_alias_name,omitempty"`
	RoomVersion               string                              `json:"room_version,omitempty"`
	Topic                     string                              `json:"topic,omitempty"`
	Visibility                string                              `json:"visibility,omitempty"`
}

type StateEvent struct {
	Content  any    `json:"content"`
	StateKey string `json:"state_key,omitempty"`
	Type     string `json:"type"`
}

type Invite3pid struct {
	Address       string `json:"address"`
	IDAccessToken string `json:"id_access_token"`
	IDServer      string `json:"id_server"`
	Medium        string `json:"medium"`
}

type RequestDirectoryRoomPut struct {
	RoomID string `json:"room_id"`
}

type RequestRooms struct {
	Reason string `json:"reason,omitempty"`
	UserID string `json:"user_id"`
}
type RequestRoomsInvite RequestRooms
type RequestRoomsKick RequestRooms
type RequestRoomsBan RequestRooms
type RequestRoomsUnban RequestRooms

type RequestJoin struct {
	Reason           string           `json:"reason,omitempty"`
	ThirdPartySigned ThirdPartySigned `json:"third_party_signed,omitempty"`
}
type RequestRoomsJoin RequestJoin

type ThirdPartySigned struct {
	Mxid       string                       `json:"mxid"`
	Sender     string                       `json:"sender"`
	Signatures map[string]map[string]string `json:"signatures"`
	Token      string                       `json:"token"`
}

type RequestKnock struct {
	Reason string `json:"reason,omitempty"`
}

type RequestRoomsLeave struct {
	Reason string `json:"reason,omitempty"`
}

type RequestDirectoryListRoomPut struct {
	Visibility string `json:"visibility"`
}

type RequestPublicRoomsPost struct {
	Filter               Filter `json:"filter,omitempty"`
	IncludeAllNetworks   bool   `json:"include_all_networks,omitempty"`
	Limit                int    `json:"limit,omitempty"`
	Since                string `json:"since,omitempty"`
	ThirdPartyInstanceID string `json:"third_party_instance_id,omitempty"`
}

type Filter struct {
	GenericSearchTerm string   `json:"generic_search_term,omitempty"`
	RoomTypes         []string `json:"room_types,omitempty"`
}

type RequestSync struct {
	Filter      string    `json:"filter,omitempty"`
	FullState   bool      `json:"full_state,omitempty"`
	SetPresence string    `json:"set_presence,omitempty"`
	Since       time.Time `json:"since,omitempty"`
	TimeOut     int       `json:"time_out,omitempty"`
}

type RequestGetMessage struct {
	Dir    string    `json:"dir,omitempty"`
	Filter string    `json:"filter,omitempty"`
	From   time.Time `json:"from,omitempty"`
	Limit  int       `json:"limit,omitempty"`
	To     time.Time `json:"to,omitempty"`
}

type RequestSendState map[string]interface{}

type RequestSendMessage struct {
	Body    string `json:"body"`
	MsgType string `json:"msgtype"`
}
