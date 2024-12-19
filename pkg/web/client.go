package web

import (
	"encoding/json"
	"net/http"

	"hotaru.hana.im/server/pkg/web/model"
)

func clientVersionsHandler(w http.ResponseWriter, _ *http.Request) {
	resp := &model.ResponseClientVersions{
		Versions: []string{"v1.11"},
	}
	raw, _ := json.Marshal(resp)
	w.Write(raw)
}
