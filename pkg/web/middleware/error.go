package middleware

import (
	"encoding/json"
	"net/http"

	"hotaru.hana.im/server/pkg/web/model"
)

func ErrorTooLarge(w http.ResponseWriter) {
	resp := &model.ResponseError{ErrCode: model.ErrCodeTooLarge}
	raw, _ := json.Marshal(resp)
	http.Error(w, string(raw), http.StatusRequestEntityTooLarge)
}
