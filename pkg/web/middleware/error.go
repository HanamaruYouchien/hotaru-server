package middleware

import (
	"net/http"

	"hotaru.hana.im/server/pkg/web/model"
)

func Error(w http.ResponseWriter, status int, errCode string, msg string) {
	w.WriteHeader(status)
	resp := &model.ResponseError{ErrCode: errCode, Error: msg}
	RenderJSON(w, resp)
}

func ErrorTooLarge(w http.ResponseWriter) {
	Error(w, http.StatusRequestEntityTooLarge, model.ErrCodeTooLarge, "")
}

func ErrorForbiddenMsg(w http.ResponseWriter, msg string) {
	Error(w, http.StatusForbidden, model.ErrCodeForbidden, msg)
}

func ErrorForbidden(w http.ResponseWriter) {
	ErrorForbiddenMsg(w, "")
}

func ErrorNotJson(w http.ResponseWriter) {
	Error(w, http.StatusBadRequest, model.ErrCodeNotJson, "")
}

func ErrorBadJson(w http.ResponseWriter) {
	Error(w, http.StatusBadRequest, model.ErrCodeBadJson, "")
}

func ErrorUserInUse(w http.ResponseWriter) {
	Error(w, http.StatusBadRequest, model.ErrCodeUserInUse, "username is already taken")
}

func ErrorUnknownMsg(w http.ResponseWriter, msg string) {
	Error(w, http.StatusBadRequest, model.ErrCodeUnknown, msg)
}

func ErrorUnknown(w http.ResponseWriter) {
	ErrorUnknownMsg(w, "")
}
