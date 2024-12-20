package model

const (
	AuthenticationTypePassword           = "m.login.password"
	AuthenticationTypeRecaptcha          = "m.login.recaptcha"
	AuthenticationTypeSSO                = "m.login.sso"
	AuthenticationTypeIdentity           = "m.login.email.identity"
	AuthenticationTypeMsisdn             = "m.login.msisdn"
	AuthenticationTypeDummy              = "m.login.dummy"
	AuthenticationTypeRegistrationToken  = "m.login.registration_token"
	AuthenticationTypeApplicationService = "m.login.application_service"

	IdentifierTypeUser       = "m.id.user"
	IdentifierTypeThirdParty = "m.id.thirdparty"
	IdentifierTypePhone      = "m.id.phone"
)

type RequestInteractiveAuthentication struct {
	Type          string          `json:"type"`
	Session       string          `json:"session"`
	Identifier    Identifier      `json:"identifier,omitempty"`
	ThreepidCreds Request3pidBind `json:"threepid_creds,omitempty"`
	Password      string          `json:"password,omitempty"`
	Response      string          `json:"response,omitempty"`
	Token         string          `json:"token,omitempty"`
}

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

type AuthenticationData struct {
	Session string `json:"session"`
	Type    string `json:"type"`
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
