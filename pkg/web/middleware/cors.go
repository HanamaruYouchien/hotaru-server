package middleware

import "net/http"

func WithCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization")

		next.ServeHTTP(w, r)
	})
}
