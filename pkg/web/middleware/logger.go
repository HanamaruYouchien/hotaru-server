package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

const CtxKeyLogger = "logger"

func Logger(l *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t := time.Now()
			defer func() {
				l.
					Info().
					Int("status", ww.Status()).
					Str("method", r.Method).
					Str("url", r.URL.Path).
					Dur("request_time", time.Since(t)).
					Int("body_bytes_sent", ww.BytesWritten()).
					Send()
			}()

			next.ServeHTTP(ww, WithLogEntry(r, l))
		})
	}
}

func WithLogEntry(r *http.Request, l *zerolog.Logger) *http.Request {
	r = r.WithContext(context.WithValue(r.Context(), CtxKeyLogger, l))
	return r
}
