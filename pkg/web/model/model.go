package model

type ResponseError struct {
	ErrCode string `json:"errcode"`
	Error   string `json:"error"`
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
	ErrCodeInUse                       = "M_USER_IN_USE"
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
)

type ResponseClientVersions struct {
	Versions []string `json:"versions"`
}
