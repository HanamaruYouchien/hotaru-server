package web

import (
	"net/http"

	"hotaru.hana.im/server/pkg/web/middleware"
	"hotaru.hana.im/server/pkg/web/model"
)

func clientVersionsHandler(w http.ResponseWriter, _ *http.Request) {
	resp := &model.ResponseClientVersions{
		Versions: []string{"v1.11"},
	}
	middleware.RenderJSON(w, resp)
}
