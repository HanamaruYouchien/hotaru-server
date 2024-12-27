package model

import (
	"encoding/json"

	"hotaru.hana.im/server/pkg/storage"
)

type ResponseError struct {
	ErrCode string `json:"errcode"`
	Error   string `json:"error,omitempty"`
}

const (
	ErrCodeForbidden     = "M_FORBIDDEN"
	ErrCodeUnknownToken  = "M_UNKNOWN_TOKEN"
	ErrCodeMissingToken  = "M_MISSING_TOKEN"
	ErrCodeUserLocked    = "M_USER_LOCKED"
	ErrCodeBadJson       = "M_BAD_JSON"
	ErrCodeNotJson       = "M_NOT_JSON"
	ErrCodeNotFound      = "M_NOT_FOUND"
	ErrCodeLimitExceeded = "M_LIMIT_EXCEEDED"
	ErrCodeUnrecognized  = "M_UNRECOGNIZED"
	ErrCodeUnknown       = "M_UNKNOWN"

	ErrCodeUnauthorized                = "M_UNAUTHORIZED"
	ErrCodeUserDeactivated             = "M_USER_DEACTIVATED"
	ErrCodeUserInUse                   = "M_USER_IN_USE"
	ErrCodeRoomState                   = "M_INVALID_ROOM_STATE"
	ErrCodeThreepidInUse               = "M_THREEPID_IN_USE"
	ErrCodeThreepidNotFound            = "M_THREEPID_NOT_FOUND"
	ErrCodeThreepidAuthFailed          = "M_THREEPID_AUTH_FAILED"
	ErrCodeThreepidDenied              = "M_THREEPID_DENIED"
	ErrCodeServerNotTrusted            = "M_SERVER_NOT_TRUSTED"
	ErrCodeUnsupportedRoomVersion      = "M_UNSUPPORTED_ROOM_VERSION"
	ErrCodeImcompatibalRoomVersion     = "M_INCOMPATIBLE_ROOM_VERSION"
	ErrCodeBadState                    = "M_BAD_STATE"
	ErrCodeGuestAccessForbidden        = "M_GUEST_ACCESS_FORBIDDEN"
	ErrCodeCaptchaNeeded               = "M_CAPTCHA_NEEDED"
	ErrCodeCaptchaInvalid              = "M_CAPTCHA_INVALID"
	ErrCodeMissingParam                = "M_MISSING_PARAM"
	ErrCodeInvalidParam                = "M_INVALID_PARAM"
	ErrCodeTooLarge                    = "M_TOO_LARGE"
	ErrCodeExclusive                   = "M_EXCLUSIVE"
	ErrCodeResourceLimitExeeded        = "M_RESOURCE_LIMIT_EXCEEDED"
	ErrCodeCannotLeaveServerNoticeRoom = "M_CANNOT_LEAVE_SERVER_NOTICE_ROOM"
	ErrCodeWeakPassword                = "M_WEAK_PASSWORD"
	ErrCodeInvalidUsername             = "M_INVALID_USERNAME"
	ErrCodeRoomInUse                   = "M_ROOM_IN_USE"
)

type ResponseClientVersions struct {
	Versions []string `json:"versions"`
}

type ResponseLoginGet struct {
	Flows []LoginFlow `json:"flows"`
}

type LoginFlow struct {
	Type          string `json:"type"`
	GetLoginToken bool   `json:"get_login_token,omitempty"`
}

type ResponseLoginPost struct {
	AccessToken  string              `json:"access_token"`
	DeviceID     string              `json:"device_id"`
	ExpiresInMs  int                 `json:"expires_in_ms,omitempty"`
	RefreshToken string              `json:"refresh_token,omitempty"`
	UserID       string              `json:"user_id"`
	WellKnown    DiscoveryInfomation `json:"well_known,omitempty"`
}

type DiscoveryInfomation struct {
	Homeserver     HomeserverInfomation      `json:"m.homeserver"`
	IdentityServer IdentityServerInformation `json:"m.identity_server,omitempty"`
}

type HomeserverInfomation struct {
	BaseURL string `json:"base_url"`
}

type IdentityServerInformation struct {
	BaseURL string `json:"base_url"`
}

type ResponseGetToken struct {
	ExpiresInMs int    `json:"expires_in_ms"`
	LoginToken  string `json:"login_token"`
}

// Interactive Auth API
// ResponseGetTokenUnauthorized, ResponseDeactivateUnauthorized, ResponsePasswordUnauthorized, ResponseRegisterUnauthorized, Response3pidAddUnauthorized
type ResponseUnauthorized struct {
	Completed []string          `json:"completed,omitempty"`
	Flows     []FlowInformation `json:"flows"`
	Params    map[string]any    `json:"params,omitempty"`
	Session   string            `json:"session,omitempty"`
}

type FlowInformation struct {
	Stages []string `json:"stages"`
}

type ResponseRefresh struct {
	AccessToken  string `json:"access_token"`
	ExpiresInMs  int    `json:"expires_in_ms,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// ResponseDeactivate, Response3pidDelete, Response3pidUnbind
type ResponseDeactivate struct {
	IDDerverUnbindResult string `json:"id_server_unbind_result"`
}
type Response3pidDelete ResponseDeactivate
type Response3pidUnbind RequestDeactivate

// ResponsePasswordEmailRequestToken, ResponsePasswordMsisdnRequestToken, ResponseRegisterEmailRequestToken, ResponseRegisterMsisdnRequestToken, Response3pidEmailRequestToken, ResponseMsisdnEmailRequestToken
type ResponseRequestToken struct {
	Sid       string `json:"sid"`
	SubmitURL string `json:"submit_url,omitempty"`
}

type ResponseRegister struct {
	AccessToken  string `json:"access_token,omitempty"`
	DeviceID     string `json:"device_id,omitempty"`
	ExpiresInMs  int    `json:"expires_in_ms,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	UserID       string `json:"user_id"`
}

type ResponseRegisterAvailable struct {
	Available bool `json:"available,omitempty"`
}

type Response3pid struct {
	AddedAt     int    `json:"added_at"`
	Address     string `json:"address"`
	Medium      string `json:"medium"`
	ValidatedAt int    `json:"validated_at"`
}

type ResponseWhoami struct {
	DeviceID string `json:"device_id,omitempty"`
	IsGuest  bool   `json:"is_guest,omitempty"`
	UserID   string `json:"user_id"`
}

type ResponseProfile struct {
	AvatarURL   string `json:"avatar_url,omitempty"`
	DisplayName string `json:"displayname,omitempty"`
}

type ResponseUserDirectorySearch struct {
	Limited bool   `json:"limited"`
	Results []User `json:"results"`
}

type User struct {
	UserID      string `json:"user_id"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	DisplayName string `json:"displayname,omitempty"`
}

type ResponseRoom struct {
	RoomID string `json:"room_id"`
}
type ResponseCreateRoom ResponseRoom
type ResponseJoinRoom ResponseRoom
type ResponseRoomsJoin ResponseRoom
type ResponseKnock ResponseRoom

type ResponseDirectoryRoomGet struct {
	RoomID  string   `json:"room_id"`
	Servers []string `json:"servers"`
}

type ResponseRoomsAliases struct {
	Aliases []string `json:"aliases"`
}

type ResponseJoinedRooms struct {
	JoinedRooms []string `json:"joined_rooms"`
}

type ResponseDirectoryListRoomGet struct {
	Visibility string `json:"visibility"`
}

type ResponsePublicRooms struct {
	Chunk                  []PublicRoomsChunk `json:"chunk"`
	NextBatch              string             `json:"next_batch"`
	PrevBatch              string             `json:"prev_batch"`
	TotalRoomCountEstimate int                `json:"total_room_count_estimate"`
}
type ResponsePublicRoomsGet ResponsePublicRooms
type ResponsePublicRoomsPost ResponsePublicRooms

type PublicRoomsChunk struct {
	AvatarURL        string `json:"avatar_url,omitempty"`
	CanonicalAlias   string `json:"canonical_alias,omitempty"`
	GuestCanJoin     bool   `json:"guest_can_join"`
	JoinRule         string `json:"join_rule,omitempty"`
	Name             string `json:"name,omitempty"`
	NumJoinedMembers int    `json:"num_joined_members"`
	RoomID           string `json:"room_id"`
	RoomType         string `json:"room_type,omitempty"`
	Topic            string `json:"topic,omitempty"`
	WorldReadable    bool   `json:"world_readable"`
}

type BooleanCapability struct {
	Enabled bool `json:"enabled"`
}

type RoomVersionsCapability struct {
	Available []string `json:"available"`
	Default   string   `json:"default"`
}

type CapabilitiesResponse struct {
	Capabilities struct {
		ThreePIDChanges BooleanCapability      `json:"m.3pid_changes"`
		ChangePassword  BooleanCapability      `json:"m.change_password"`
		GetLoginToken   BooleanCapability      `json:"m.get_login_token"`
		RoomVersions    RoomVersionsCapability `json:"m.room_versions"`
		SetAvatarURL    BooleanCapability      `json:"m.set_avatar_url"`
		SetDisplayName  BooleanCapability      `json:"m.set_displayname"`
	} `json:"capabilities"`
}

type EventResponse struct {
	Event storage.Event `json:"events"`
}

type SyncResponse struct {
	AccountData AccountData `json:"account_data,omitempty"` // The global private data created by this user.
	NextBatch   string      `json:"next_batch"`
	Presence    Presence    `json:"presence,omitempty"` // The updates to the presence status of other users.
	Rooms       Rooms       `json:"rooms,omitempty"`

	// The Next two are used to e2ee
	// DeviceLists DeviceLists `json:"device_lists,omitempty"`
	// DeviceOneTimeKeysCount map[string]int `json:"device_one_time_keys_count,omitempty"`

	// Send-To-Device required
	// ToDevice ToDevice `json:"to_device,omitempty"`
}

type Rooms struct {
	Invite map[string]InviteState `json:"invite,omitempty"`
	Join   map[string]JoinedRoom  `json:"join,omitempty"`
	Knock  map[string]KnockedRoom `json:"knock,omitempty"`
	Leave  map[string]LeftRoom    `json:"leave,omitempty"`
}

type InvitedRoom struct {
	InviteState InviteState `json:"invite_state,omitempty"`
}

type JoinedRoom struct {
	AccountData               AccountData                         `json:"account_data,omitempty"`
	Ephemeral                 Ephemeral                           `json:"ephemeral,omitempty"`
	State                     State                               `json:"state,omitempty"`
	Summary                   RoomSummary                         `json:"summary,omitempty"`
	Timeline                  Timeline                            `json:"timeline,omitempty"`
	UnreadNotifications       UnreadNotificationCounts            `json:"unread_notifications,omitempty"`
	UnreadThreadNotifications map[string]ThreadNotificationCounts `json:"unread_thread_notifications,omitempty"`
}

type KnockedRoom struct {
	KnockState KnockState `json:"knock_state,omitempty"`
}

type LeftRoom struct {
	AccountData AccountData `json:"account_data,omitempty"`
	State       State       `json:"state,omitempty"`
	Timeline    Timeline    `json:"timeline,omitempty"`
}

type InviteState struct {
	Events []StrippedStateEvent `json:"events,omitempty"`
}

type KnockState struct {
	Events []StrippedStateEvent `json:"events,omitempty"`
}

type StrippedStateEvent struct {
	Content  json.RawMessage `json:"content"`
	Sender   string          `json:"sender"`
	StateKey string          `json:"state_key"`
	Type     string          `json:"type"`
}

type Timeline struct {
	Events    []storage.Event `json:"events"` // no room id
	Limited   bool            `json:"limited,omitempty"`
	PrevBatch string          `json:"prev_batch,omitempty"`
}

type State struct {
	Events []storage.Event `json:"events,omitempty"` // no room id
}

type RoomSummary struct {
	Heroes             []string `json:"m.heroes,omitempty"`
	InvitedMemberCount int      `json:"m.invited_member_count,omitempty"`
	JoinedMemberCount  int      `json:"m.joined_member_count,omitempty"`
}

type AccountData struct {
	Events []Event `json:"events,omitempty"`
}

type Presence struct {
	Events []Event `json:"events,omitempty"`
}

type Ephemeral struct {
	Events []Event `json:"events,omitempty"`
}

type Event struct {
	// This type is correspond to `Event` in https://spec.matrix.org/v1.13/client-server-api/#get_matrixclientv3sync,
	// with only content and type in a standard ClientEvent.
	// It might can be replaced by storage.Event for lazy work..
	Content json.RawMessage `json:"content"`
	Type    string          `json:"type"`
}

type UnreadNotificationCounts struct {
	HighlightCount    int `json:"highlight_count,omitempty"`
	NotificationCount int `json:"notification_count,omitempty"`
}

type ThreadNotificationCounts struct {
	HighlightCount    int `json:"highlight_count,omitempty"`
	NotificationCount int `json:"notification_count,omitempty"`
}

type Member struct {
	AvatarURL   *string `json:"avatar_url"`
	DisplayName string  `json:"display_name"`
}

type JoinedMembersResponse struct {
	Joined map[string]Member `json:"joined"`
}

type MembersResponse struct {
	Events []storage.Event `json:"chunk"`
}

type MessagesResponse struct {
	Events []storage.Event `json:"chunk"`
	Start  string          `json:"start"`
	End    string          `json:"end"`
}

type EventSentResponse struct {
	EventID string `json:"event_id"`
}
