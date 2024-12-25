package storage

import "time"

type Event struct {
	EventId      string `xorm:"pk"`
	RoomId       string
	Type         string
	Sender       string
	StateKey     string
	CreatedAt    time.Time `xorm:"created"`
	Content      any       `xorm:"text"`
	UnsignedData any       `xorm:"text"`
}

const (
	EventTypeRoomCanonicalAlias = "m.room.canonical_alias"
	EventTypeRoomCreate         = "m.room.create"
	EventTypeRoomJoinRules      = "m.room.join_rules"
	EventTypeRoomMember         = "m.room.member"
	EventTypeRoomPowerLevels    = "m.room.power_levels"

	EventTypeRoomMessage      = "m.room.message"
	EventTypeRoomName         = "m.room.name"
	EventTypeRoomTopic        = "m.room.topic"
	EventTypeRoomAvatar       = "m.room.avatar"
	EventTypeRoomPinnedEvents = "m.room.pinned_events"

	MessageTypeText     = "m.text"
	MessageTypeEmote    = "m.emote"
	MessageTypeNotice   = "m.notice"
	MessageTypeImage    = "m.image"
	MessageTypeFile     = "m.file"
	MessageTypeAudio    = "m.audio"
	MessageTypeLocation = "m.location"
	MessageTypeVideo    = "m.video"
)

type EventRoomCanonicalAliasContent struct {
	Alias      string   `json:"alias,omitempty"`
	AltAliases []string `json:"alt_aliases,omitempty"`
}

type EventRoomCreateContent struct {
	Creator     string   `json:"creator"`
	Federate    bool     `json:"m.federate,omitempty"`
	RoomVersion string   `json:"room_version,omitempty"`
	Type        string   `json:"type,omitempty"`
	Predecessor struct { // PreviousRoom
		RoomID  string `json:"room_id"`
		EventID string `json:"event_id"`
	} `json:"predecessor,omitempty"`
}

type EventRoomJoinRulesContent struct {
	Allow []struct { // AllowCondition
		RoomID string `json:"room_id,omitempty"`
		Type   string `json:"type"`
	} `json:"allow,omitempty"`
	JoinRule string `json:"join_rule"`
}

type EventRoomMemberContent struct {
	AvatarURL                    string   `json:"avatar_url,omitempty"`
	DisplayName                  string   `json:"displayname,omitempty"`
	IsDirect                     bool     `json:"is_direct,omitempty"`
	JoinAuthorisedViaUsersServer string   `json:"join_authorised_via_users_server,omitempty"`
	Membership                   string   `json:"membership"`
	Reason                       string   `json:"reason,omitempty"`
	ThirdPartyInvite             struct { // Invite
		DisplayName string   `json:"display_name"`
		Signed      struct { // signed
			Mxid       string                       `json:"mxid"`
			Signatures map[string]map[string]string `json:"signatures"`
			Token      string                       `json:"token"`
		}
	} `json:"third_party_invite,omitempty"`
}

type EventRoomPowerLevelsContent struct {
	Ban           int            `json:"ban,omitempty"`
	Events        map[string]int `json:"events,omitempty"`
	EventsDefault int            `json:"events_default,omitempty"`
	Invite        int            `json:"invite,omitempty"`
	Kick          int            `json:"kick,omitempty"`
	Redact        int            `json:"redact,omitempty"`
	StateDefault  int            `json:"state_default,omitempty"`
	Users         map[string]int `json:"users,omitempty"`
	UsersDefault  int            `json:"users_default,omitempty"`
	Notifications struct {       // Notifications
		Room int `json:"room,omitempty"`
		// TODO: Other properties
	} `json:"notifications,omitempty"`
}

type EventRoomMessageContent struct {
	Body    string `json:"body"`
	Msgtype string `json:"msgtype"`
}

type EventRoomMessageTextContent struct {
	EventRoomMessageContent

	Format        string `json:"format,omitempty"`
	FormattedBody string `json:"formatted_body,omitempty"`
}

type EventRoomMessageEmoteContent struct {
	EventRoomMessageContent

	Format        string `json:"format,omitempty"`
	FormattedBody string `json:"formatted_body,omitempty"`
}

type EventRoomMessageNoticeContent struct {
	EventRoomMessageContent

	Format        string `json:"format,omitempty"`
	FormattedBody string `json:"formatted_body,omitempty"`
}

type EventRoomMessageImageContent struct {
	EventRoomMessageContent

	Format        string `json:"format,omitempty"`
	FormattedBody string `json:"formatted_body,omitempty"`
	Filename      string `json:"filename,omitempty"`
	URL           string `json:"url,omitempty"`
	// File EncryptedFile `json:"EncryptedFile,omitempty"`
	Info struct { // ImageInfo
		H             int           `json:"h,omitempty"`
		W             int           `json:"w,omitempty"`
		MimeType      string        `json:"mimetype,omitempty"`
		Size          int           `json:"size,omitempty"`
		ThumbnailURL  string        `json:"thumbnail_url,omitempty"`
		ThumbnailInfo ThumbnailInfo `json:"thumbnail_info,omitempty"`
		// ThumbnailFile EncryptedFile `json:"thumbnail_file,omitempty"`
	} `json:"info,omitempty"`
}

type EventRoomMessageFileContent struct {
	EventRoomMessageContent

	Format        string `json:"format,omitempty"`
	FormattedBody string `json:"formatted_body,omitempty"`
	Filename      string `json:"filename,omitempty"`
	URL           string `json:"url,omitempty"`
	// File EncryptedFile `json:"EncryptedFile,omitempty"`
	Info struct { // FileInfo
		MimeType      string        `json:"mimetype,omitempty"`
		Size          int           `json:"size,omitempty"`
		ThumbnailURL  string        `json:"thumbnail_url,omitempty"`
		ThumbnailInfo ThumbnailInfo `json:"thumbnail_info,omitempty"`
		// ThumbnailFile EncryptedFile `json:"thumbnail_file,omitempty"`
	} `json:"info,omitempty"`
}

type EventRoomMessageAudioContent struct {
	EventRoomMessageContent

	Format        string `json:"format,omitempty"`
	FormattedBody string `json:"formatted_body,omitempty"`
	Filename      string `json:"filename,omitempty"`
	URL           string `json:"url,omitempty"`
	// File EncryptedFile `json:"EncryptedFile,omitempty"`
	Info struct { // AudioInfo
		Duration int    `json:"duration,omitempty"`
		MimeType string `json:"mimetype,omitempty"`
		Size     int    `json:"size,omitempty"`
	} `json:"info"`
}

type EventRoomMessageLocationContent struct {
	EventRoomMessageContent

	GeoURI string   `json:"geo_uri"`
	Info   struct { // LocationInfo
		ThumbnailURL  string        `json:"thumbnail_url,omitempty"`
		ThumbnailInfo ThumbnailInfo `json:"thumbnail_info,omitempty"`
		// ThumbnailFile EncryptedFile `json:"thumbnail_file,omitempty"`
	} `json:"info,omitempty"`
}

type EventRoomMessageVideoContent struct {
	EventRoomMessageContent

	Format        string `json:"format,omitempty"`
	FormattedBody string `json:"formatted_body,omitempty"`
	Filename      string `json:"filename,omitempty"`
	URL           string `json:"url,omitempty"`
	// File EncryptedFile `json:"EncryptedFile,omitempty"`
	Info struct { // VideoInfo
		H             int           `json:"h,omitempty"`
		W             int           `json:"w,omitempty"`
		Duration      int           `json:"duration,omitempty"`
		MimeType      string        `json:"mimetype,omitempty"`
		Size          int           `json:"size,omitempty"`
		ThumbnailURL  string        `json:"thumbnail_url,omitempty"`
		ThumbnailInfo ThumbnailInfo `json:"thumbnail_info,omitempty"`
		// ThumbnailFile EncryptedFile `json:"thumbnail_file,omitempty"`
	} `json:"info,omitempty"`
}

type EventRoomNameContent struct {
	Name string `json:"name"`
}

type EventRoomTopicContent struct {
	Topic string `json:"topic"`
}

type EventRoomAvatarContent struct {
	URL  string   `json:"url,omitempty"`
	Info struct { // AvatarInfo
		H             int           `json:"h,omitempty"`
		W             int           `json:"w,omitempty"`
		MimeType      string        `json:"mimetype,omitempty"`
		Size          int           `json:"size,omitempty"`
		ThumbnailURL  string        `json:"thumbnail_url,omitempty"`
		ThumbnailInfo ThumbnailInfo `json:"thumbnail_info,omitempty"`
	} `json:"info,omitempty"`
}

type EventRoomPinnedEventsContent struct {
	Pinned []string `json:"pinned"`
}

type ThumbnailInfo struct {
	H        int    `json:"h,omitempty"`
	W        int    `json:"w,omitempty"`
	MimeType string `json:"mimetype,omitempty"`
	Size     int    `json:"size,omitempty"`
}
