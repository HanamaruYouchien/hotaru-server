package middleware

import (
	"encoding/json"
	"net/http"

	"hotaru.hana.im/server/pkg/web/model"
)

func ErrorTooLarge(w http.ResponseWriter) {
	resp := &model.ResponseError{ErrCode: model.ErrCodeTooLarge}
	raw, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusRequestEntityTooLarge)
	w.Write(raw)
}

func ErrorNotJson(w http.ResponseWriter) {
	resp := &model.ResponseError{ErrCode: model.ErrCodeNotJson}
	raw, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(raw)
}

func ErrorBadJson(w http.ResponseWriter) {
	resp := &model.ResponseError{ErrCode: model.ErrCodeBadJson}
	raw, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(raw)
}

func ErrorUserInUse(w http.ResponseWriter) {
	resp := &model.ResponseError{ErrCode: model.ErrCodeUserInUse}
	raw, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(raw)
}

func ErrorUnknown(w http.ResponseWriter) {
	resp := &model.ResponseError{ErrCode: model.ErrCodeUnknown}
	raw, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(raw)
}
