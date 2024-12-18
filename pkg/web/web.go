package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Serve() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/_hotaru", hotaruRouter())

	http.ListenAndServe(":8008", r)
}

func hotaruRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	return r
}
